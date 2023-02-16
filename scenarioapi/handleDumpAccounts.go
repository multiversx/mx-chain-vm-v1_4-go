package main

import (
	oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"
)

func (ae *Executor) HandleDumpAccounts() (oj.OJsonObject, error) {
	ojAccounts := oj.OJsonList{}
	for _, account := range ae.world.AcctMap {
		ojAccounts = append(ojAccounts, worldAccountToOJ(account))
	}
	return &ojAccounts, nil
}
