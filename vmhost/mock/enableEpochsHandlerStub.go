package mock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
)

var _ vmhost.EnableEpochsHandler = (*EnableEpochsHandlerStub)(nil)

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	IsFlagEnabledInCurrentEpochCalled              func(flag core.EnableEpochFlag) bool
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

// IsFlagEnabledInCurrentEpoch -
func (stub *EnableEpochsHandlerStub) IsFlagEnabledInCurrentEpoch(flag core.EnableEpochFlag) bool {
	if stub.IsFlagEnabledInCurrentEpochCalled != nil {
		return stub.IsFlagEnabledInCurrentEpochCalled(flag)
	}
	return false
}

// GetActivationEpoch -
func (stub *EnableEpochsHandlerStub) GetActivationEpoch(flag core.EnableEpochFlag) uint32 {
	if stub.GetActivationEpochCalled != nil {
		return stub.GetActivationEpochCalled(flag)
	}
	switch flag {
	case core.MultiESDTTransferFixOnCallBackFlag:
		return stub.MultiESDTTransferAsyncCallBackEnableEpochField
	case core.FixOOGReturnCodeFlag:
		return stub.FixOOGReturnCodeEnableEpochField
	case core.RemoveNonUpdatedStorageFlag:
		return stub.RemoveNonUpdatedStorageEnableEpochField
	case core.CreateNFTThroughExecByCallerFlag:
		return stub.CreateNFTThroughExecByCallerEnableEpochField
	case core.FailExecutionOnEveryAPIErrorFlag:
		return stub.FixFailExecutionOnErrorEnableEpochField
	case core.ManagedCryptoAPIsFlag:
		return stub.ManagedCryptoAPIEnableEpochField
	case core.DisableExecByCallerFlag:
		return stub.DisableExecByCallerEnableEpochField
	case core.RefactorContextFlag:
		return stub.RefactorContextEnableEpochField
	case core.CheckExecuteOnReadOnlyFlag:
		return stub.CheckExecuteReadOnlyEnableEpochField
	case core.StorageAPICostOptimizationFlag:
		return stub.StorageAPICostOptimizationEnableEpochField
	default:
		return 0
	}
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
