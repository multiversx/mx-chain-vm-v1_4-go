package arwenpart

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// ArwenPart is
type ArwenPart struct {
	Messenger *ChildMessenger
	VMHost    vmcommon.VMExecutionHandler
}

// NewArwenPart creates
func NewArwenPart(input *os.File, output *os.File) (*ArwenPart, error) {
	messenger := NewChildMessenger(input, output)
	blockchain := NewBlockchainHookGateway(messenger)
	crypto := NewCryptoHookGateway()
	arwenVirtualMachineType := []byte{5, 0} // TODO
	blockGasLimit := uint64(10000000)       // TODO
	gasSchedule := config.MakeGasMap(1)     // TODO

	host, err := host.NewArwenVM(blockchain, crypto, arwenVirtualMachineType, blockGasLimit, gasSchedule)
	if err != nil {
		return nil, err
	}

	return &ArwenPart{
		Messenger: messenger,
		VMHost:    host,
	}, nil
}

// StartLoop runs the main loop
func (part *ArwenPart) StartLoop() error {
	err := part.doLoop()
	part.Messenger.Shutdown()
	common.LogError("Arwen: end of loop, err=%v", err)
	return err
}

// doLoop ends only when a critical failure takes place
func (part *ArwenPart) doLoop() error {
	for {
		request, err := part.Messenger.ReceiveContractRequest()
		if err != nil {
			return err
		}

		response, err := part.handleContractRequest(request)
		if err != nil {
			return err
		}

		// Successful execution, send response
		part.Messenger.SendContractResponse(response)
		part.Messenger.EndDialogue()
	}
}

func (part *ArwenPart) handleContractRequest(request *common.ContractRequest) (*common.HookCallRequestOrContractResponse, error) {
	common.LogDebug("Arwen: handleContractRequest() %v", request)

	switch request.Action {
	case "Deploy":
		return part.doRunSmartContractCreate(request), nil
	case "Call":
		return part.doRunSmartContractCall(request), nil
	case "Stop":
		return nil, common.ErrStopPerNodeRequest
	default:
		return nil, common.ErrBadRequestFromNode
	}
}

func (part *ArwenPart) doRunSmartContractCreate(request *common.ContractRequest) *common.HookCallRequestOrContractResponse {
	vmOutput, err := part.VMHost.RunSmartContractCreate(request.CreateInput)
	common.LogDebug("doRunSmartContractCreate, err=%v", err)
	common.LogDebugJSON("VMOutput", vmOutput)
	return common.NewContractResponse(vmOutput, err)
}

func (part *ArwenPart) doRunSmartContractCall(request *common.ContractRequest) *common.HookCallRequestOrContractResponse {
	vmOutput, err := part.VMHost.RunSmartContractCall(request.CallInput)
	common.LogDebug("doRunSmartContractCall, err=%v", err)
	return common.NewContractResponse(vmOutput, err)
}
