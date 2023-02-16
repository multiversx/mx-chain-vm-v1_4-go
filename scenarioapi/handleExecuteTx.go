package main

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost"
)

type Tx struct {
	DisplayLogs   bool
	Type      		TxType
	From      		[]byte
	To        		[]byte
	EGLDAmount 	  *big.Int
	ESDTs         []*TxESDT
	Code      		[]byte
	FunctionName  string
	FunctionArgs  [][]byte
	GasPrice  		uint64
	GasLimit  		uint64
}

type TxESDT struct {
	Id 		  []byte
	Nonce   uint64
	Amount  *big.Int
}

type TxType int

const (
	DeploySc TxType = iota
	CallSc
	Transfer
)

func (ae *Executor) HandleExecuteTx(reqBody []byte) (oj.OJsonObject, error) {
	tx, err := parseTx(reqBody)
	if err != nil {
		return nil, err
	}
	return ae.executeTx(tx)
}

func parseTx(reqBody []byte) (*Tx, error) {
	raw, err := oj.ParseOrderedJSON(reqBody)
	if err != nil {
		return nil, err
	}
	stepMap, isMap := raw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block info object is not a map")
	}
	tx := &Tx{
		EGLDAmount: big.NewInt(0),
	}
	for _, kvp := range stepMap.OrderedKV {
		switch kvp.Key {
			case "displayLogs":
				tx.DisplayLogs, err = ojBoolToBool(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction displayLogs: %w", err)
				}
			case "type":
				typeStr, err := ojStrToStr(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction type: %w", err)
				}
				if typeStr == "callSc" {
					tx.Type = CallSc
				} else if typeStr == "deploySc" {
					tx.Type = DeploySc
				} else if typeStr == "transfer" {
					tx.Type = Transfer
				} else {
					return nil, fmt.Errorf("invalid transaction type: %w", err)
				}
			case "from":
				tx.From, err = ojHexStrToBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction from: %w", err)
				}
			case "to":
				tx.To, err = ojHexStrToBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction to: %w", err)
				}
			case "egldAmount":
				tx.EGLDAmount, err = ojIntStrToBigint(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction egldAmount: %w", err)
				}
			case "esdts":
				list, isList := kvp.Value.(*oj.OJsonList);
				if !isList {
					return nil, errors.New("invalid transaction esdts")
				}
				esdts := []*TxESDT{}
				for _, esdtRaw := range list.AsList() {
					esdt, err := parseTxESDT(esdtRaw)
					if err != nil {
						return nil, fmt.Errorf("invalid transaction esdt: %w", err)
					}
					esdts = append(esdts, esdt)
				}
				tx.ESDTs = esdts
				if err != nil {
					return nil, fmt.Errorf("invalid transaction esdts: %w", err)
				}
			case "code":
				tx.Code, err = ojHexStrToBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction code: %w", err)
				}
			case "functionName":
				tx.FunctionName, err = ojStrToStr(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction functionName: %w", err)
				}
			case "functionArgs":
				list, isList := kvp.Value.(*oj.OJsonList);
				if !isList {
					return nil, errors.New("invalid transaction functionArguments")
				}
				var arguments [][]byte
				for _, item := range list.AsList() {
					argument, err := ojHexStrToBytes(item)
					if err != nil {
						return nil, fmt.Errorf("invalid transaction functionArgument: %w", err)
					}
					arguments = append(arguments, argument)
				}
				tx.FunctionArgs = arguments
				if err != nil {
					return nil, fmt.Errorf("invalid transaction functionArguments: %w", err)
				}
			case "gasLimit":
				tx.GasLimit, err = ojIntStrToUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction gasLimit: %w", err)
				}
			case "gasPrice":
				tx.GasPrice, err = ojIntStrToUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid transaction gasPrice: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown transaction field: %s", kvp.Key)
		}
	}
	return tx, nil
}

func parseTxESDT(esdtRaw oj.OJsonObject) (*TxESDT, error) {
	esdtMap, isMap := esdtRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, fmt.Errorf("wrong ESDT Multi-Transfer format")
	}
	esdtData := TxESDT{}
	var err error
	for _, kvp := range esdtMap.OrderedKV {
		switch kvp.Key {
			case "id":
				esdtData.Id, err = ojStrToBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid ESDT token name: %w", err)
				}
			case "nonce":
				esdtData.Nonce, err = ojIntStrToUint64(kvp.Value)
				if err != nil {
					return nil, errors.New("invalid account nonce")
				}
			case "amount":
				esdtData.Amount, err = ojIntStrToBigint(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid ESDT balance: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown transaction ESDT data field: %s", kvp.Key)
		}
	}
	return &esdtData, nil
}


func (ae *Executor) executeTx(tx *Tx) (oj.OJsonObject, error) {
	if tx.DisplayLogs {
		vmhost.SetLoggingForTests()
	}
	vmOutput, err := ae._executeTx(tx)
	if tx.DisplayLogs {
		vmhost.DisableLoggingForTests()
	}
	if err != nil {
		return nil, err
	}
	return vmOutputToOJ(vmOutput), err
}

func (ae *Executor) _executeTx(tx *Tx) (*vmcommon.VMOutput, error) {
	ae.world.CreateStateBackup()
	var err error
	defer func() {
		if err != nil {
			errRollback := ae.world.RollbackChanges()
			if errRollback != nil {
				err = errRollback
			}
		} else {
			errCommit := ae.world.CommitChanges()
			if errCommit != nil {
				err = errCommit
			}
		}
	}()
	sender := ae.world.AcctMap.GetAccount(tx.From)
	switch tx.Type {
		case DeploySc:
			newAddressMock := &worldmock.NewAddressMock{
				CreatorAddress: sender.Address,
				CreatorNonce:   sender.Nonce,
				NewAddress:     tx.To,
			}
			ae.world.NewAddressMocks = append(ae.world.NewAddressMocks, newAddressMock)
	}
	err = ae.world.UpdateWorldStateBefore(tx.From, tx.GasLimit, tx.GasPrice)
	if err != nil {
		return nil, fmt.Errorf("could not set up tx: %w", err)
	}
	gasForExecution := tx.GasLimit
	if tx.ESDTs != nil {
		vmInputArguments := [][]byte{
			tx.To,
			big.NewInt(0).SetUint64(uint64(len(tx.ESDTs))).Bytes(),
		}
		for _, esdtTransfer := range tx.ESDTs {
			vmInputArguments = append(
				vmInputArguments,
				esdtTransfer.Id,
				big.NewInt(0).SetUint64(esdtTransfer.Nonce).Bytes(),
				esdtTransfer.Amount.Bytes(),
			)
		}
		vmCallInput := &vmcommon.ContractCallInput{
			VMInput: vmcommon.VMInput{
				CallerAddr:  tx.From,
				Arguments:   vmInputArguments,
				CallValue:   big.NewInt(0),
				CallType:    vm.DirectCall,
				GasPrice:    tx.GasPrice,
				GasProvided: gasForExecution,
				GasLocked:   0,
			},
			RecipientAddr:     tx.From,
			Function:          core.BuiltInFunctionMultiESDTNFTTransfer,
			AllowInitFunction: false,
		}
		vmOutput, err := ae.world.BuiltinFuncs.ProcessBuiltInFunction(vmCallInput)
		if err != nil {
			return nil, err
		}
		if vmOutput.ReturnCode != vmcommon.Ok {
			return nil, fmt.Errorf(
				"MultiESDTtransfer failed: retcode = %d, msg = %s",
				vmOutput.ReturnCode, vmOutput.ReturnMessage)
		}
		gasForExecution = vmOutput.GasRemaining
	}
	// we use fake vm outputs for transactions that don't use the VM, just for convenience
	var vmOutput *vmcommon.VMOutput
	if sender.Balance.Cmp(tx.EGLDAmount) < 0 {
		// out of funds is handled by the protocol, so it needs to be mocked here
		vmOutput = &vmcommon.VMOutput{
			ReturnData:      make([][]byte, 0),
			ReturnCode:      vmcommon.OutOfFunds,
			ReturnMessage:   "",
			GasRemaining:    0,
			GasRefund:       big.NewInt(0),
			OutputAccounts:  make(map[string]*vmcommon.OutputAccount),
			DeletedAccounts: make([][]byte, 0),
			TouchedAccounts: make([][]byte, 0),
			Logs:            make([]*vmcommon.LogEntry, 0),
		}
	} else {
		switch tx.Type {
			case DeploySc:
				txHash := make([]byte, 0)
				vmCallInput := &vmcommon.ContractCreateInput{
					ContractCode: tx.Code,
					VMInput:      vmcommon.VMInput{
						CallerAddr:     tx.From,
						Arguments:      tx.FunctionArgs,
						CallValue:      tx.EGLDAmount,
						GasPrice:       tx.GasPrice,
						GasProvided:    gasForExecution,
						OriginalTxHash: txHash,
						CurrentTxHash:  txHash,
						ESDTTransfers:  esdtTransfersToVmEsdtTransfers(tx.ESDTs),
					},
				}
				vmOutput, err = ae.vmHost.RunSmartContractCreate(vmCallInput)
				if err != nil {
					return nil, err
				}
			case CallSc:
				txHash := make([]byte, 0)
				vmCallInput := &vmcommon.ContractCallInput{
					RecipientAddr: tx.To,
					Function:      tx.FunctionName,
					VMInput:       vmcommon.VMInput{
						CallerAddr:     tx.From,
						Arguments:      tx.FunctionArgs,
						CallValue:      tx.EGLDAmount,
						GasPrice:       tx.GasPrice,
						GasProvided:    gasForExecution,
						OriginalTxHash: txHash,
						CurrentTxHash:  txHash,
						ESDTTransfers:  esdtTransfersToVmEsdtTransfers(tx.ESDTs),
					},
				}
				vmOutput, err = ae.vmHost.RunSmartContractCall(vmCallInput)
				if err != nil {
					return nil, err
				}
			case Transfer:
				outputAccounts := make(map[string]*vmcommon.OutputAccount)
				outputAccounts[string(tx.To)] = &vmcommon.OutputAccount{
					Address:      tx.To,
					BalanceDelta: tx.EGLDAmount,
				}
				vmOutput = &vmcommon.VMOutput{
					ReturnData:      make([][]byte, 0),
					ReturnCode:      vmcommon.Ok,
					ReturnMessage:   "",
					GasRemaining:    0,
					GasRefund:       big.NewInt(0),
					OutputAccounts:  outputAccounts,
					DeletedAccounts: make([][]byte, 0),
					TouchedAccounts: make([][]byte, 0),
					Logs:            make([]*vmcommon.LogEntry, 0),
				}
			default:
				return nil, errors.New("unknown transaction type")
		}
	}
	if vmOutput.ReturnCode == vmcommon.Ok {
		err := ae.world.UpdateBalanceWithDelta(tx.From, big.NewInt(0).Neg(tx.EGLDAmount))
		if err != nil {
			return nil, err
		}
		err = ae.world.UpdateAccounts(vmOutput.OutputAccounts, vmOutput.DeletedAccounts)
		if err != nil {
			return nil, err
		}
	} else {
		err = fmt.Errorf(
			"tx step failed: retcode=%d, msg=%s",
			vmOutput.ReturnCode, vmOutput.ReturnMessage)
	}
	return vmOutput, nil
}

func esdtTransfersToVmEsdtTransfers(esdtTransfers []*TxESDT) []*vmcommon.ESDTTransfer {
	vmESDTTransfers := make([]*vmcommon.ESDTTransfer, len(esdtTransfers))
	for i := 0; i < len(esdtTransfers); i++ {
		vmESDTTransfers[i] = &vmcommon.ESDTTransfer{}
		vmESDTTransfers[i].ESDTTokenName = esdtTransfers[i].Id
		vmESDTTransfers[i].ESDTTokenNonce = esdtTransfers[i].Nonce
		vmESDTTransfers[i].ESDTValue = esdtTransfers[i].Amount
		if vmESDTTransfers[i].ESDTTokenNonce != 0 {
			vmESDTTransfers[i].ESDTTokenType = uint32(core.NonFungible)
		} else {
			vmESDTTransfers[i].ESDTTokenType = uint32(core.Fungible)
		}
	}
	return vmESDTTransfers
}

func vmOutputToOJ(vmOutput *vmcommon.VMOutput) oj.OJsonObject {
	ojOutput := oj.NewMap()
	ojOutput.Put("code", intToOjIntStr(int(vmOutput.ReturnCode)))
	ojOutput.Put("message", strToOjStr(vmOutput.ReturnMessage))
	var returnData oj.OJsonList
	for _, data := range vmOutput.ReturnData {
		returnData = append(returnData, bytesToOjHexStr(data))
	}
	ojOutput.Put("data", &returnData)
	return ojOutput
}
