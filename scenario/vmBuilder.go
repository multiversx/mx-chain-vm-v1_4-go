package scenario

import (
	"fmt"

	"github.com/multiversx/mx-chain-vm-v1_4-go/config"
	gasSchedules "github.com/multiversx/mx-chain-vm-v1_4-go/scenario/gasSchedules"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/hostCore"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/mock"

	"github.com/multiversx/mx-chain-core-go/core"
	scenexec "github.com/multiversx/mx-chain-scenario-go/executor"
	mj "github.com/multiversx/mx-chain-scenario-go/model"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
)

var _ scenexec.VMBuilder = (*ScenarioVMHostBuilder)(nil)

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// VMTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with VM.
type ScenarioVMHostBuilder struct {
}

// NewScenarioVMHostBuilder creates a default ScenarioVMHostBuilder.
func NewScenarioVMHostBuilder() *ScenarioVMHostBuilder {
	return &ScenarioVMHostBuilder{}
}

// NewMockWorld defines how the MockWorld is initialized.
func (*ScenarioVMHostBuilder) NewMockWorld() *worldmock.MockWorld {
	return mock.NewMockWorldVM14()
}

// GasScheduleMapFromScenarios provides the correct gas schedule for the gas schedule named specified in a scenario.
func (svb *ScenarioVMHostBuilder) GasScheduleMapFromScenarios(scenGasSchedule mj.GasSchedule) (worldmock.GasScheduleMap, error) {
	switch scenGasSchedule {
	case mj.GasScheduleDefault:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	case mj.GasScheduleDummy:
		return config.MakeGasMapForTests(), nil
	case mj.GasScheduleV3:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV3())
	case mj.GasScheduleV4:
		return gasSchedules.LoadGasScheduleConfig(gasSchedules.GetV4())
	default:
		return nil, fmt.Errorf("unknown scenario GasSchedule: %d", scenGasSchedule)
	}
}

// NewVM will create a new VM instance with pointers to a mock world and given gas schedule.
func (svb *ScenarioVMHostBuilder) NewVM(
	world *worldmock.MockWorld,
	gasSchedule map[string]map[string]uint64,
) (scenexec.VMInterface, error) {
	blockGasLimit := uint64(10000000)
	esdtTransferParser, _ := parsers.NewESDTTransferParser(worldmock.WorldMarshalizer)

	return hostCore.NewVMHost(
		world,
		&vmhost.VMHostParameters{
			VMType:                   TestVMType,
			BlockGasLimit:            blockGasLimit,
			GasSchedule:              gasSchedule,
			BuiltInFuncContainer:     world.BuiltinFuncs.Container,
			ProtectedKeyPrefix:       []byte(core.ProtectedKeyPrefix),
			ESDTTransferParser:       esdtTransferParser,
			EpochNotifier:            &mock.EpochNotifierStub{},
			EnableEpochsHandler:      world.EnableEpochsHandler,
			WasmerSIGSEGVPassthrough: false,
			Hasher:                   worldmock.DefaultHasher,
		})

}

// DefaultScenarioExecutor provides a scenario executor with VM 1.5, default configuration
func DefaultScenarioExecutor() *scenexec.ScenarioExecutor {
	return scenexec.NewScenarioExecutor(NewScenarioVMHostBuilder())
}
