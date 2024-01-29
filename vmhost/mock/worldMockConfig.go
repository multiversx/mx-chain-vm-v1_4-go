package mock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
)

// ChangeUsernameFlag is the DNS v1 disable epoch
const ChangeUsernameFlag core.EnableEpochFlag = "ChangeUsernameFlag"

// NewMockWorldVM14 creates a MockWorld specifically configured for all VM 1.4 tests.
func NewMockWorldVM14() *worldmock.MockWorld {
	world := worldmock.NewMockWorld()
	world.EnableEpochsHandler = &EnableEpochsHandlerStub{
		IsFlagDefinedCalled: func(flag core.EnableEpochFlag) bool {
			return true
		},
		IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
			return flag != ChangeUsernameFlag // relevant in DNS test
		},
	}
	return world
}
