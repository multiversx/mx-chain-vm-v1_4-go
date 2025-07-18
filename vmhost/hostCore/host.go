package hostCore

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
	"github.com/multiversx/mx-chain-vm-v1_4-go/crypto"
	"github.com/multiversx/mx-chain-vm-v1_4-go/crypto/factory"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/contexts"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooksmeta"
	"github.com/multiversx/mx-chain-vm-v1_4-go/wasmer"
)

var log = logger.GetOrCreate("vm/host")
var logGasTrace = logger.GetOrCreate("gasTrace")

// MaximumWasmerInstanceCount represents the maximum number of Wasmer instances that can be active at the same time
var MaximumWasmerInstanceCount = uint64(10)

var _ vmhost.VMHost = (*vmHost)(nil)

const minExecutionTimeout = time.Second
const internalVMErrors = "internalVMErrors"

// allFlags must have all flags used by mx-chain-vm-v1_4-go in the current version
var allFlags = []core.EnableEpochFlag{
	vmhost.MultiESDTTransferFixOnCallBackFlag,
	vmhost.RemoveNonUpdatedStorageFlag,
	vmhost.CreateNFTThroughExecByCallerFlag,
	vmhost.StorageAPICostOptimizationFlag,
	vmhost.CheckExecuteOnReadOnlyFlag,
	vmhost.FailExecutionOnEveryAPIErrorFlag,
	vmhost.ManagedCryptoAPIsFlag,
	vmhost.DisableExecByCallerFlag,
	vmhost.RefactorContextFlag,
	vmhost.RuntimeMemStoreLimitFlag,
	vmhost.RuntimeCodeSizeFixFlag,
	vmhost.FixOOGReturnCodeFlag,
}

// vmHost implements HostContext interface.
type vmHost struct {
	cryptoHook       crypto.VMCrypto
	mutExecution     sync.RWMutex
	closingInstance  bool
	executionTimeout time.Duration

	ethInput []byte

	blockchainContext   vmhost.BlockchainContext
	runtimeContext      vmhost.RuntimeContext
	outputContext       vmhost.OutputContext
	meteringContext     vmhost.MeteringContext
	storageContext      vmhost.StorageContext
	managedTypesContext vmhost.ManagedTypesContext

	gasSchedule          config.GasScheduleMap
	scAPIMethods         *wasmer.Imports
	builtInFuncContainer vmcommon.BuiltInFunctionContainer
	esdtTransferParser   vmcommon.ESDTTransferParser
	enableEpochsHandler  vmhost.EnableEpochsHandler
	activationEpochMap   map[uint32]struct{}
}

// NewVMHost creates a new VM vmHost
func NewVMHost(
	blockChainHook vmcommon.LegacyBlockchainHook,
	hostParameters *vmhost.VMHostParameters,
) (vmhost.VMHost, error) {

	if check.IfNil(blockChainHook) {
		return nil, vmhost.ErrNilBlockChainHook
	}
	if hostParameters == nil {
		return nil, vmhost.ErrNilHostParameters
	}
	if check.IfNil(hostParameters.ESDTTransferParser) {
		return nil, vmhost.ErrNilESDTTransferParser
	}
	if check.IfNil(hostParameters.BuiltInFuncContainer) {
		return nil, vmhost.ErrNilBuiltInFunctionsContainer
	}
	if check.IfNil(hostParameters.EpochNotifier) {
		return nil, vmhost.ErrNilEpochNotifier
	}
	if check.IfNil(hostParameters.EnableEpochsHandler) {
		return nil, vmhost.ErrNilEnableEpochsHandler
	}
	err := core.CheckHandlerCompatibility(hostParameters.EnableEpochsHandler, allFlags)
	if err != nil {
		return nil, err
	}
	if check.IfNil(hostParameters.Hasher) {
		return nil, vmhost.ErrNilHasher
	}
	if hostParameters.VMType == nil {
		return nil, vmhost.ErrNilVMType
	}

	cryptoHook := factory.NewVMCrypto()
	host := &vmHost{
		cryptoHook:           cryptoHook,
		meteringContext:      nil,
		runtimeContext:       nil,
		blockchainContext:    nil,
		storageContext:       nil,
		managedTypesContext:  nil,
		gasSchedule:          hostParameters.GasSchedule,
		scAPIMethods:         nil,
		builtInFuncContainer: hostParameters.BuiltInFuncContainer,
		esdtTransferParser:   hostParameters.ESDTTransferParser,
		executionTimeout:     minExecutionTimeout,
		enableEpochsHandler:  hostParameters.EnableEpochsHandler,
	}

	host.activationEpochMap = createActivationMap(hostParameters)

	newExecutionTimeout := time.Duration(hostParameters.TimeOutForSCExecutionInMilliseconds) * time.Millisecond
	if newExecutionTimeout > minExecutionTimeout {
		host.executionTimeout = newExecutionTimeout
	}

	imports := vmhooksmeta.NewEIFunctions()
	err = PopulateAllImports(imports)
	if err != nil {
		return nil, err
	}

	wasmerImports := wasmer.ConvertImports(imports)
	err = wasmer.SetImports(wasmerImports)
	if err != nil {
		return nil, err
	}

	host.scAPIMethods = wasmerImports

	host.blockchainContext, err = contexts.NewBlockchainContext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.runtimeContext, err = contexts.NewRuntimeContext(
		host,
		hostParameters.VMType,
		host.builtInFuncContainer,
		hostParameters.Hasher,
	)
	if err != nil {
		return nil, err
	}

	host.meteringContext, err = contexts.NewMeteringContext(host, hostParameters.GasSchedule, hostParameters.BlockGasLimit)
	if err != nil {
		return nil, err
	}

	host.outputContext, err = contexts.NewOutputContext(host)
	if err != nil {
		return nil, err
	}

	host.storageContext, err = contexts.NewStorageContext(
		host,
		blockChainHook,
		hostParameters.ProtectedKeyPrefix,
	)
	if err != nil {
		return nil, err
	}

	host.managedTypesContext, err = contexts.NewManagedTypesContext(host)
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(host.gasSchedule)
	if err != nil {
		return nil, err
	}

	host.runtimeContext.SetMaxInstanceCount(MaximumWasmerInstanceCount)

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)
	wasmer.SetRkyvSerializationEnabled(true)

	if hostParameters.WasmerSIGSEGVPassthrough {
		wasmer.SetSIGSEGVPassthrough()
	}

	wasmer.ForceInstallSighandlers()

	host.initContexts()
	hostParameters.EpochNotifier.RegisterNotifyHandler(host)

	return host, nil
}

func createActivationMap(hostParameters *vmhost.VMHostParameters) map[uint32]struct{} {
	activationMap := make(map[uint32]struct{})

	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.CheckExecuteOnReadOnlyFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.DisableExecByCallerFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.RefactorContextFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.FailExecutionOnEveryAPIErrorFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.ManagedCryptoAPIsFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.CreateNFTThroughExecByCallerFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.FixOOGReturnCodeFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.MultiESDTTransferFixOnCallBackFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.RemoveNonUpdatedStorageFlag)] = struct{}{}
	activationMap[hostParameters.EnableEpochsHandler.GetActivationEpoch(vmhost.StorageAPICostOptimizationFlag)] = struct{}{}

	return activationMap
}

// GetVersion returns the VM version string
func (host *vmHost) GetVersion() string {
	return vmhost.VMVersion
}

// Crypto returns the VMCrypto instance of the host
func (host *vmHost) Crypto() crypto.VMCrypto {
	return host.cryptoHook
}

// Blockchain returns the BlockchainContext instance of the host
func (host *vmHost) Blockchain() vmhost.BlockchainContext {
	return host.blockchainContext
}

// Runtime returns the RuntimeContext instance of the host
func (host *vmHost) Runtime() vmhost.RuntimeContext {
	return host.runtimeContext
}

// Output returns the OutputContext instance of the host
func (host *vmHost) Output() vmhost.OutputContext {
	return host.outputContext
}

// Metering returns the MeteringContext instance of the host
func (host *vmHost) Metering() vmhost.MeteringContext {
	return host.meteringContext
}

// Storage returns the StorageContext instance of the host
func (host *vmHost) Storage() vmhost.StorageContext {
	return host.storageContext
}

// EnableEpochsHandler returns the enableEpochsHandler instance of the host
func (host *vmHost) EnableEpochsHandler() vmhost.EnableEpochsHandler {
	return host.enableEpochsHandler
}

// ManagedTypes returns the ManagedTypeContext instance of the host
func (host *vmHost) ManagedTypes() vmhost.ManagedTypesContext {
	return host.managedTypesContext
}

// GetContexts returns the main contexts of the host
func (host *vmHost) GetContexts() (
	vmhost.ManagedTypesContext,
	vmhost.BlockchainContext,
	vmhost.MeteringContext,
	vmhost.OutputContext,
	vmhost.RuntimeContext,
	vmhost.StorageContext,
) {
	return host.managedTypesContext,
		host.blockchainContext,
		host.meteringContext,
		host.outputContext,
		host.runtimeContext,
		host.storageContext
}

// InitState resets the contexts of the host and reconfigures its flags
func (host *vmHost) InitState() {
	host.initContexts()
}

func (host *vmHost) close() {
	host.runtimeContext.ClearWarmInstanceCache()
}

// Close will close all underlying processes
func (host *vmHost) Close() error {
	host.mutExecution.Lock()
	host.close()
	host.closingInstance = true
	host.mutExecution.Unlock()

	return nil
}

// Reset is a function which closes the VM and resets the closingInstance variable
func (host *vmHost) Reset() {
	host.mutExecution.Lock()
	host.close()
	// keep closingInstance flag to false
	host.mutExecution.Unlock()
}

func (host *vmHost) initContexts() {
	host.ClearContextStateStack()
	host.managedTypesContext.InitState()
	host.outputContext.InitState()
	host.meteringContext.InitState()
	host.runtimeContext.InitState()
	host.storageContext.InitState()
	host.ethInput = nil
}

// ClearContextStateStack cleans the state stacks of all the contexts of the host
func (host *vmHost) ClearContextStateStack() {
	host.managedTypesContext.ClearStateStack()
	host.outputContext.ClearStateStack()
	host.meteringContext.ClearStateStack()
	host.runtimeContext.ClearStateStack()
	host.storageContext.ClearStateStack()
}

// GetAPIMethods returns the EEI as a set of imports for Wasmer
func (host *vmHost) GetAPIMethods() *wasmer.Imports {
	return host.scAPIMethods
}

// GasScheduleChange applies a new gas schedule to the host
func (host *vmHost) GasScheduleChange(newGasSchedule config.GasScheduleMap) {
	host.mutExecution.Lock()
	defer host.mutExecution.Unlock()

	host.gasSchedule = newGasSchedule
	gasCostConfig, err := config.CreateGasConfig(newGasSchedule)
	if err != nil {
		log.Error("cannot apply new gas config", "err", err)
		return
	}

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host.meteringContext.SetGasSchedule(newGasSchedule)
	host.runtimeContext.ClearWarmInstanceCache()
}

// GetGasScheduleMap returns the currently stored gas schedule
func (host *vmHost) GetGasScheduleMap() config.GasScheduleMap {
	return host.gasSchedule
}

// GetGasScheduleMap returns the curent gas trace, used in scenario tests
func (host *vmHost) GetGasTrace() map[string]map[string][]uint64 {
	return host.meteringContext.GetGasTrace()
}

// SetGasTracing configures the gas tracing flag, used in scenario tests
func (host *vmHost) SetGasTracing(enableGasTracing bool) {
	host.meteringContext.SetGasTracing(enableGasTracing)
}

// RunSmartContractCreate executes the deployment of a new contract
func (host *vmHost) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	host.mutExecution.RLock()
	defer host.mutExecution.RUnlock()

	if host.closingInstance {
		return nil, vmhost.ErrVMIsClosing
	}

	host.setGasTracerEnabledIfLogIsTrace()
	ctx, cancel := context.WithTimeout(context.Background(), host.executionTimeout)
	defer cancel()

	log.Trace("RunSmartContractCreate begin",
		"len(code)", len(input.ContractCode),
		"metadata", input.ContractCodeMetadata,
		"gasProvided", input.GasProvided,
		"gasLocked", input.GasLocked)

	done := make(chan struct{})
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				log.Error("VM execution panicked", "error", r, "stack", "\n"+string(debug.Stack()))
				err = vmhost.ErrExecutionPanicked
				host.Runtime().CleanInstance()
			} else {
				host.Runtime().EndExecution()
			}

			close(done)
		}()

		vmOutput = host.doRunSmartContractCreate(input)
		logsFromErrors := host.createLogEntryFromErrors(input.CallerAddr, input.CallerAddr, "_init")
		if logsFromErrors != nil {
			vmOutput.Logs = append(vmOutput.Logs, logsFromErrors)
		}

		log.Trace("RunSmartContractCreate end",
			"returnCode", vmOutput.ReturnCode,
			"returnMessage", vmOutput.ReturnMessage,
			"gasRemaining", vmOutput.GasRemaining)

		host.logFromGasTracer("init")
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		host.Runtime().FailExecution(vmhost.ErrExecutionFailedWithTimeout)
		<-done
		err = vmhost.ErrExecutionFailedWithTimeout
	}

	return
}

// RunSmartContractCall executes the call of an existing contract
func (host *vmHost) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	host.mutExecution.RLock()
	defer host.mutExecution.RUnlock()

	if host.closingInstance {
		return nil, vmhost.ErrVMIsClosing
	}

	host.setGasTracerEnabledIfLogIsTrace()
	ctx, cancel := context.WithTimeout(context.Background(), host.executionTimeout)
	defer cancel()

	log.Trace("RunSmartContractCall begin",
		"function", input.Function,
		"gasProvided", input.GasProvided,
		"gasLocked", input.GasLocked)

	done := make(chan struct{})
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				log.Error("VM execution panicked", "error", r, "stack", "\n"+string(debug.Stack()))
				err = vmhost.ErrExecutionPanicked
				host.Runtime().CleanInstance()
			} else {
				host.Runtime().EndExecution()
			}

			close(done)
		}()

		isUpgrade := input.Function == vmhost.UpgradeFunctionName
		if isUpgrade {
			vmOutput = host.doRunSmartContractUpgrade(input)
		} else {
			vmOutput = host.doRunSmartContractCall(input)
		}

		logsFromErrors := host.createLogEntryFromErrors(input.CallerAddr, input.RecipientAddr, input.Function)
		if logsFromErrors != nil {
			vmOutput.Logs = append(vmOutput.Logs, logsFromErrors)
		}

		log.Trace("RunSmartContractCall end",
			"function", input.Function,
			"returnCode", vmOutput.ReturnCode,
			"returnMessage", vmOutput.ReturnMessage,
			"gasRemaining", vmOutput.GasRemaining)

		host.logFromGasTracer(input.Function)
	}()

	select {
	case <-done:
		// Normal termination.
		return
	case <-ctx.Done():
		// Terminated due to timeout. The VM sets the `ExecutionFailed` breakpoint
		// in Wasmer. Also, the VM must wait for Wasmer to reach the end of a WASM
		// basic block in order to close the WASM instance cleanly. This is done by
		// reading the `done` channel once more, awaiting the call to `close(done)`
		// from above.
		host.Runtime().FailExecution(vmhost.ErrExecutionFailedWithTimeout)
		<-done
		err = vmhost.ErrExecutionFailedWithTimeout
	}

	return
}

func (host *vmHost) createLogEntryFromErrors(sndAddress, rcvAddress []byte, function string) *vmcommon.LogEntry {
	formattedErrors := host.runtimeContext.GetAllErrors()
	if formattedErrors == nil {
		return nil
	}

	logFromError := &vmcommon.LogEntry{
		Identifier: []byte(internalVMErrors),
		Address:    sndAddress,
		Topics:     [][]byte{rcvAddress, []byte(function)},
		Data:       [][]byte{[]byte(formattedErrors.Error())},
	}

	return logFromError
}

// AreInSameShard returns true if the provided addresses are part of the same shard
func (host *vmHost) AreInSameShard(leftAddress []byte, rightAddress []byte) bool {
	blockchain := host.Blockchain()
	leftShard := blockchain.GetShardOfAddress(leftAddress)
	rightShard := blockchain.GetShardOfAddress(rightAddress)

	return leftShard == rightShard
}

// IsInterfaceNil returns true if there is no value under the interface
func (host *vmHost) IsInterfaceNil() bool {
	return host == nil
}

// SetRuntimeContext sets the runtimeContext for this host, used in tests
func (host *vmHost) SetRuntimeContext(runtime vmhost.RuntimeContext) {
	host.runtimeContext = runtime
}

// GetRuntimeErrors obtains the cumultated error object after running the SC
func (host *vmHost) GetRuntimeErrors() error {
	if host.runtimeContext != nil {
		return host.runtimeContext.GetAllErrors()
	}
	return nil
}

// SetBuiltInFunctionsContainer sets the built in function container - only for testing
func (host *vmHost) SetBuiltInFunctionsContainer(builtInFuncs vmcommon.BuiltInFunctionContainer) {
	if check.IfNil(builtInFuncs) {
		return
	}
	host.builtInFuncContainer = builtInFuncs
}

// EpochConfirmed is called whenever a new epoch is confirmed
func (host *vmHost) EpochConfirmed(epoch uint32, _ uint64) {
	_, ok := host.activationEpochMap[epoch]
	if ok {
		host.Runtime().ClearWarmInstanceCache()
		host.Blockchain().ClearCompiledCodes()
	}
}

// FixOOGReturnCodeEnabled returns true if the corresponding flag is set
func (host *vmHost) FixOOGReturnCodeEnabled() bool {
	return host.enableEpochsHandler.IsFlagEnabled(vmhost.FixOOGReturnCodeFlag)
}

// FixFailExecutionEnabled returns true if the corresponding flag is set
func (host *vmHost) FixFailExecutionEnabled() bool {
	return host.enableEpochsHandler.IsFlagEnabled(vmhost.FailExecutionOnEveryAPIErrorFlag)
}

// CreateNFTOnExecByCallerEnabled returns true if the corresponding flag is set
func (host *vmHost) CreateNFTOnExecByCallerEnabled() bool {
	return host.enableEpochsHandler.IsFlagEnabled(vmhost.CreateNFTThroughExecByCallerFlag)
}

// DisableExecByCaller returns true if the corresponding flag is set
func (host *vmHost) DisableExecByCaller() bool {
	return host.enableEpochsHandler.IsFlagEnabled(vmhost.DisableExecByCallerFlag)
}

// CheckExecuteReadOnly returns true if the corresponding flag is set
func (host *vmHost) CheckExecuteReadOnly() bool {
	return host.enableEpochsHandler.IsFlagEnabled(vmhost.CheckExecuteOnReadOnlyFlag)
}

func (host *vmHost) setGasTracerEnabledIfLogIsTrace() {
	host.Metering().SetGasTracing(false)
	if logGasTrace.GetLevel() == logger.LogTrace {
		host.Metering().SetGasTracing(true)
	}
}

func (host *vmHost) logFromGasTracer(functionName string) {
	if logGasTrace.GetLevel() == logger.LogTrace {
		scGasTrace := host.meteringContext.GetGasTrace()
		totalGasUsedByAPIs := 0
		for scAddress, gasTrace := range scGasTrace {
			logGasTrace.Trace("Gas Trace for", "SC Address", scAddress, "function", functionName)
			for apiName, value := range gasTrace {
				totalGasUsed := uint64(0)
				for _, usedGas := range value {
					totalGasUsed += usedGas
				}
				logGasTrace.Trace("Gas Trace for", "apiName", apiName, "totalGasUsed", totalGasUsed, "numberOfCalls", len(value))
				totalGasUsedByAPIs += int(totalGasUsed)
			}
			logGasTrace.Trace("Gas Trace for", "TotalGasUsedByAPIs", totalGasUsedByAPIs)
		}
	}
}
