package main

import oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"

func (ae *Executor) HandleReset() (oj.OJsonObject, error) {
	ae.world.Clear()
	ae.vmHost.Reset()
	ojBool := oj.OJsonBool(true)
	return &ojBool, nil
}
