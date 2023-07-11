package mock

import (
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
)

var _ vmhost.EnableEpochsHandler = (*EnableEpochsHandlerStub)(nil)

// EnableEpochsHandlerStub -
type EnableEpochsHandlerStub struct {
	IsCheckCorrectTokenIDForTransferRoleFlagEnabledField bool
	IsMultiESDTTransferFixOnCallBackFlagEnabledField     bool
	IsFixOOGReturnCodeFlagEnabledField                   bool
	IsRemoveNonUpdatedStorageFlagEnabledField            bool
	IsCreateNFTThroughExecByCallerFlagEnabledField       bool
	IsStorageAPICostOptimizationFlagEnabledField         bool
	IsFailExecutionOnEveryAPIErrorFlagEnabledField       bool
	IsManagedCryptoAPIsFlagEnabledField                  bool
	IsDisableExecByCallerFlagEnabledField                bool
	IsRefactorContextFlagEnabledField                    bool
	IsCheckExecuteOnReadOnlyFlagEnabledField             bool
	IsRuntimeMemStoreLimitEnabledField                   bool
	IsRuntimeCodeSizeFixEnabledField                     bool
	IsAutoBalanceDataTriesEnabledField                   bool
	GetCurrentEpochField                                 uint32
	MultiESDTTransferAsyncCallBackEnableEpochField       uint32
	FixOOGReturnCodeEnableEpochField                     uint32
	RemoveNonUpdatedStorageEnableEpochField              uint32
	CreateNFTThroughExecByCallerEnableEpochField         uint32
	FixFailExecutionOnErrorEnableEpochField              uint32
	ManagedCryptoAPIEnableEpochField                     uint32
	DisableExecByCallerEnableEpochField                  uint32
	RefactorContextEnableEpochField                      uint32
	CheckExecuteReadOnlyEnableEpochField                 uint32
	StorageAPICostOptimizationEnableEpochField           uint32
}

// GetCurrentEpoch -
func (stub *EnableEpochsHandlerStub) GetCurrentEpoch() uint32 {
	return stub.GetCurrentEpochField
}

// IsStorageAPICostOptimizationFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsStorageAPICostOptimizationFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsStorageAPICostOptimizationFlagEnabledField
}

// IsManagedCryptoAPIsFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsManagedCryptoAPIsFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsManagedCryptoAPIsFlagEnabledField
}

// IsMultiESDTTransferFixOnCallBackFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsMultiESDTTransferFixOnCallBackFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsMultiESDTTransferFixOnCallBackFlagEnabledField
}

// IsRemoveNonUpdatedStorageFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsRemoveNonUpdatedStorageFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsRemoveNonUpdatedStorageFlagEnabledField
}

// IsRefactorContextFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsRefactorContextFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsRefactorContextFlagEnabledField
}

// IsFailExecutionOnEveryAPIErrorFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsFailExecutionOnEveryAPIErrorFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsFailExecutionOnEveryAPIErrorFlagEnabledField
}

// IsFixOOGReturnCodeFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsFixOOGReturnCodeFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsFixOOGReturnCodeFlagEnabledField
}

// IsCreateNFTThroughExecByCallerFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsCreateNFTThroughExecByCallerFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsCreateNFTThroughExecByCallerFlagEnabledField
}

// IsDisableExecByCallerFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsDisableExecByCallerFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsDisableExecByCallerFlagEnabledField
}

// IsCheckExecuteOnReadOnlyFlagEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsCheckExecuteOnReadOnlyFlagEnabledInEpoch(_ uint32) bool {
	return stub.IsCheckExecuteOnReadOnlyFlagEnabledField
}

// IsRuntimeCodeSizeFixEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsRuntimeCodeSizeFixEnabledInEpoch(_ uint32) bool {
	return stub.IsRuntimeCodeSizeFixEnabledField
}

// IsRuntimeMemStoreLimitEnabledInEpoch -
func (stub *EnableEpochsHandlerStub) IsRuntimeMemStoreLimitEnabledInEpoch(_ uint32) bool {
	return stub.IsRuntimeMemStoreLimitEnabledField
}

// IsAutoBalanceDataTriesEnabled -
func (stub *EnableEpochsHandlerStub) IsAutoBalanceDataTriesEnabled() bool {
	return stub.IsAutoBalanceDataTriesEnabledField
}

// MultiESDTTransferAsyncCallBackEnableEpoch -
func (stub *EnableEpochsHandlerStub) MultiESDTTransferAsyncCallBackEnableEpoch() uint32 {
	return stub.MultiESDTTransferAsyncCallBackEnableEpochField
}

// FixOOGReturnCodeEnableEpoch -
func (stub *EnableEpochsHandlerStub) FixOOGReturnCodeEnableEpoch() uint32 {
	return stub.FixOOGReturnCodeEnableEpochField
}

// RemoveNonUpdatedStorageEnableEpoch -
func (stub *EnableEpochsHandlerStub) RemoveNonUpdatedStorageEnableEpoch() uint32 {
	return stub.RemoveNonUpdatedStorageEnableEpochField
}

// CreateNFTThroughExecByCallerEnableEpoch -
func (stub *EnableEpochsHandlerStub) CreateNFTThroughExecByCallerEnableEpoch() uint32 {
	return stub.CreateNFTThroughExecByCallerEnableEpochField
}

// FixFailExecutionOnErrorEnableEpoch -
func (stub *EnableEpochsHandlerStub) FixFailExecutionOnErrorEnableEpoch() uint32 {
	return stub.FixFailExecutionOnErrorEnableEpochField
}

// ManagedCryptoAPIEnableEpoch -
func (stub *EnableEpochsHandlerStub) ManagedCryptoAPIEnableEpoch() uint32 {
	return stub.ManagedCryptoAPIEnableEpochField
}

// DisableExecByCallerEnableEpoch -
func (stub *EnableEpochsHandlerStub) DisableExecByCallerEnableEpoch() uint32 {
	return stub.DisableExecByCallerEnableEpochField
}

// RefactorContextEnableEpoch -
func (stub *EnableEpochsHandlerStub) RefactorContextEnableEpoch() uint32 {
	return stub.RefactorContextEnableEpochField
}

// CheckExecuteReadOnlyEnableEpoch -
func (stub *EnableEpochsHandlerStub) CheckExecuteReadOnlyEnableEpoch() uint32 {
	return stub.CheckExecuteReadOnlyEnableEpochField
}

// StorageAPICostOptimizationEnableEpoch -
func (stub *EnableEpochsHandlerStub) StorageAPICostOptimizationEnableEpoch() uint32 {
	return stub.StorageAPICostOptimizationEnableEpochField
}

// IsInterfaceNil -
func (stub *EnableEpochsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
