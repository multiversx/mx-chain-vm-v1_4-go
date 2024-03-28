package mock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
)

var _ vmhost.EnableEpochsHandler = (*EnableEpochsHandlerStub)(nil)

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	IsFlagDefinedCalled                            func(flag core.EnableEpochFlag) bool
	IsFlagEnabledCalled                            func(flag core.EnableEpochFlag) bool
	IsFlagEnabledInEpochCalled                     func(flag core.EnableEpochFlag, epoch uint32) bool
	GetActivationEpochCalled                       func(flag core.EnableEpochFlag) uint32
	MultiESDTTransferAsyncCallBackEnableEpochField uint32
	FixOOGReturnCodeEnableEpochField               uint32
	RemoveNonUpdatedStorageEnableEpochField        uint32
	CreateNFTThroughExecByCallerEnableEpochField   uint32
	FixFailExecutionOnErrorEnableEpochField        uint32
	ManagedCryptoAPIEnableEpochField               uint32
	DisableExecByCallerEnableEpochField            uint32
	RefactorContextEnableEpochField                uint32
	CheckExecuteReadOnlyEnableEpochField           uint32
	StorageAPICostOptimizationEnableEpochField     uint32
}

// IsFlagDefined -
func (stub *EnableEpochsHandlerStub) IsFlagDefined(flag core.EnableEpochFlag) bool {
	if stub.IsFlagDefinedCalled != nil {
		return stub.IsFlagDefinedCalled(flag)
	}
	return true
}

// IsFlagEnabled -
func (stub *EnableEpochsHandlerStub) IsFlagEnabled(flag core.EnableEpochFlag) bool {
	if stub.IsFlagEnabledCalled != nil {
		return stub.IsFlagEnabledCalled(flag)
	}
	return false
}

// IsFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsFlagEnabledInEpoch(flag core.EnableEpochFlag, epoch uint32) bool {
	if stub.IsFlagEnabledInEpochCalled != nil {
		return stub.IsFlagEnabledInEpochCalled(flag, epoch)
	}
	return false
}

// GetActivationEpoch -
func (stub *EnableEpochsHandlerStub) GetActivationEpoch(flag core.EnableEpochFlag) uint32 {
	if stub.GetActivationEpochCalled != nil {
		return stub.GetActivationEpochCalled(flag)
	}
	switch flag {
	case vmhost.MultiESDTTransferFixOnCallBackFlag:
		return stub.MultiESDTTransferAsyncCallBackEnableEpochField
	case vmhost.FixOOGReturnCodeFlag:
		return stub.FixOOGReturnCodeEnableEpochField
	case vmhost.RemoveNonUpdatedStorageFlag:
		return stub.RemoveNonUpdatedStorageEnableEpochField
	case vmhost.CreateNFTThroughExecByCallerFlag:
		return stub.CreateNFTThroughExecByCallerEnableEpochField
	case vmhost.FailExecutionOnEveryAPIErrorFlag:
		return stub.FixFailExecutionOnErrorEnableEpochField
	case vmhost.ManagedCryptoAPIsFlag:
		return stub.ManagedCryptoAPIEnableEpochField
	case vmhost.DisableExecByCallerFlag:
		return stub.DisableExecByCallerEnableEpochField
	case vmhost.RefactorContextFlag:
		return stub.RefactorContextEnableEpochField
	case vmhost.CheckExecuteOnReadOnlyFlag:
		return stub.CheckExecuteReadOnlyEnableEpochField
	case vmhost.StorageAPICostOptimizationFlag:
		return stub.StorageAPICostOptimizationEnableEpochField
	default:
		return 0
	}
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
