package main

import (
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
	worldhook "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	gasSchedules "github.com/multiversx/mx-chain-vm-v1_4-go/scenarioexec/gasSchedules"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/hostCore"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/mock"
)

type Executor struct {
	world             *worldhook.MockWorld
	vmHost            vmhost.VMHost
}

func NewExecutor() (*Executor, error) {
	world := worldhook.NewMockWorld()
	gasSchedule, err := gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	if err != nil {
		return nil, err
	}
	err = world.InitBuiltinFunctions(gasSchedule)
	if err != nil {
		return nil, err
	}
	esdtTransferParser, err := parsers.NewESDTTransferParser(worldhook.WorldMarshalizer)
	if err != nil {
		return nil, err
	}
	vmHost, err := hostCore.NewVMHost(world, &vmhost.VMHostParameters{
		VMType:                   []byte{0, 0},
		BlockGasLimit:            uint64(10000000),
		GasSchedule:              gasSchedule,
		BuiltInFuncContainer:     world.BuiltinFuncs.Container,
		ProtectedKeyPrefix: 			[]byte("ELROND"),
		ESDTTransferParser:       esdtTransferParser,
		EpochNotifier: 						&mock.EpochNotifierStub{},
		EnableEpochsHandler: 			&mock.EnableEpochsHandlerStub{
			IsStorageAPICostOptimizationFlagEnabledField:         true,
			IsMultiESDTTransferFixOnCallBackFlagEnabledField:     true,
			IsFixOOGReturnCodeFlagEnabledField:                   true,
			IsRemoveNonUpdatedStorageFlagEnabledField:            true,
			IsCreateNFTThroughExecByCallerFlagEnabledField:       true,
			IsManagedCryptoAPIsFlagEnabledField:                  true,
			IsFailExecutionOnEveryAPIErrorFlagEnabledField:       true,
			IsRefactorContextFlagEnabledField:                    true,
			IsCheckCorrectTokenIDForTransferRoleFlagEnabledField: true,
			IsDisableExecByCallerFlagEnabledField:                true,
			IsESDTTransferRoleFlagEnabledField:                   true,
			IsGlobalMintBurnFlagEnabledField:                     true,
			IsTransferToMetaFlagEnabledField:                     true,
			IsCheckFrozenCollectionFlagEnabledField:              true,
			IsFixAsyncCallbackCheckFlagEnabledField:              true,
			IsESDTNFTImprovementV1FlagEnabledField:               true,
			IsSaveToSystemAccountFlagEnabledField:                true,
			IsValueLengthCheckFlagEnabledField:                   true,
			IsSCDeployFlagEnabledField:                           true,
			IsRepairCallbackFlagEnabledField:                     true,
			IsAheadOfTimeGasUsageFlagEnabledField:                true,
			IsCheckFunctionArgumentFlagEnabledField:              true,
			IsCheckExecuteOnReadOnlyFlagEnabledField:             true,
		},
		WasmerSIGSEGVPassthrough: false,
		Hasher:                   worldhook.DefaultHasher,
	})
	if err != nil {
		return nil, err
	}
	ae := Executor{
		world:  world,
		vmHost: vmHost,
	}
	return &ae, nil
}
