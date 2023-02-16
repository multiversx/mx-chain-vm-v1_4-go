package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"
)

type Account struct {
	Address         []byte
	Nonce           uint64
	Balance         *big.Int
	Storage         []*StorageKeyValuePair
	Code            []byte
	Owner           []byte
}

type StorageKeyValuePair struct {
	Key   []byte
	Value []byte
}

func (ae *Executor) HandleSetAccount(strAddress string, reqBody []byte) (oj.OJsonObject, error) {
	address, err := hex.DecodeString(strAddress)
	if err != nil {
		return nil, err
	}
	account, err := parseAccount(address, reqBody)
	if err != nil {
		return nil, err
	}
	ojBool := oj.OJsonBool(true)
	err = ae.setAccount(account)
	return &ojBool, err
}

func parseAccount(address []byte, reqBody []byte) (*Account, error) {
	raw, err := oj.ParseOrderedJSON(reqBody)
	if err != nil {
		return nil, err
	}
	accountMap, isMap := raw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled account object is not a map")
	}
	account := &Account{
		Address: 				 address,
		Nonce:           0,
		Balance:         big.NewInt(0),
	}
	for _, kvp := range accountMap.OrderedKV {
		switch kvp.Key {
			case "nonce":
				account.Nonce, err = ojIntStrToUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account nonce: %w", err)
				}
			case "egldAmount":
				account.Balance, err = ojIntStrToBigint(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account balance: %w", err)
				}
			case "storage":
				storageMap, storageOk := kvp.Value.(*oj.OJsonMap)
				if !storageOk {
					return nil, errors.New("invalid account storage")
				}
				for _, storageKvp := range storageMap.OrderedKV {
					byteKey, err := hex.DecodeString(storageKvp.Key)
					if err != nil {
						return nil, fmt.Errorf("invalid account storage key: %w", err)
					}
					byteVal, err := ojHexStrToBytes(storageKvp.Value)
					if err != nil {
						return nil, fmt.Errorf("invalid account storage value: %w", err)
					}
					stElem := StorageKeyValuePair{
						Key:   byteKey,
						Value: byteVal,
					}
					account.Storage = append(account.Storage, &stElem)
				}
			case "code":
				account.Code, err = ojHexStrToBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account code: %w", err)
				}
			case "owner":
				account.Owner, err = ojHexStrToBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account owner: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown account field: %s", kvp.Key)
		}
	}
	return account, nil
}

func (ae *Executor) setAccount(account *Account) error {
	storage := make(map[string][]byte)
	for _, kvp := range account.Storage {
		storage[string(kvp.Key)] = kvp.Value
	}
	worldAccount := &worldmock.Account{
		Address:         account.Address,
		Nonce:           account.Nonce,
		Balance:         big.NewInt(0).Set(account.Balance),
		BalanceDelta:    big.NewInt(0),
		DeveloperReward: big.NewInt(0),
		Storage:         storage,
		Code:            account.Code,
		OwnerAddress:    account.Owner,
		IsSmartContract: len(account.Code) > 0,
		CodeMetadata: (&vmcommon.CodeMetadata{
			Payable:     true,
			Upgradeable: true,
			Readable:    true,
		}).ToBytes(),
		MockWorld: ae.world,
	}
	err := worldAccount.Validate()
	if err != nil {
		return err
	}
	ae.world.AcctMap.PutAccount(worldAccount)
	return nil
}
