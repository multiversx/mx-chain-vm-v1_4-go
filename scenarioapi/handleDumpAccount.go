package main

import (
	"encoding/hex"

	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"
)

func (ae *Executor) HandleDumpAccount(strAddress string) (oj.OJsonObject, error) {
	address, err := hex.DecodeString(strAddress)
	if err != nil {
		return nil, err
	}
	account := ae.world.AcctMap.GetAccount(address)
	return worldAccountToOJ(account), nil
}

func worldAccountToOJ(account *worldmock.Account) oj.OJsonObject {
	ojAccount := oj.NewMap()
	ojAccount.Put("address", bytesToOjHexStr(account.Address))
	ojAccount.Put("nonce", intToOjIntStr(int(account.Nonce)))
	ojAccount.Put("egldAmount", bigintToOjIntStr(account.Balance))
	storageOJ := oj.NewMap()
	for k, v := range account.Storage {
		storageOJ.Put(hex.EncodeToString([]byte(k)), bytesToOjHexStr(v))
	}
	ojAccount.Put("storage", storageOJ)
	ojAccount.Put("code", bytesToOjHexStr(account.Code))
	ojAccount.Put("owner", bytesToOjHexStr(account.OwnerAddress))
	return ojAccount
}
