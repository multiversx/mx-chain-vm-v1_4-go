package mock

import (
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
)

// NewMockWorldVM14 creates a MockWorld specifically configured for all VM 1.4 tests.
func NewMockWorldVM14() *worldmock.MockWorld {
	world := worldmock.NewMockWorld()
	world.EnableEpochsHandler = &EnableEpochsHandlerStub{
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
		IsFixOldTokenLiquidityEnabledField:                   true,
		IsChangeUsernameEnabledField:                         false, // relevant in DNS test
	}
	return world
}
