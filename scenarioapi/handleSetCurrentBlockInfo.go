package main

import (
	"errors"
	"fmt"

	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
	oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"
)

type BlockInfo struct {
	Timestamp  uint64
	Nonce      uint64
	Round      uint64
	Epoch      uint64
}

func (ae *Executor) HandleSetCurrentBlockInfo(reqBody []byte) (oj.OJsonObject, error) {
	blockInfo, err := parseBlockInfo(reqBody)
	if err != nil {
		return nil, err
	}
	ae.setCurrentBlockInfo(blockInfo)
	ojBool := oj.OJsonBool(true)
	return &ojBool, nil
}

func parseBlockInfo(reqBody []byte) (*BlockInfo, error) {
	raw, err := oj.ParseOrderedJSON(reqBody)
	if err != nil {
		return nil, err
	}
	blockMap, isMap := raw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block info object is not a map")
	}
	blockInfo := &BlockInfo{}
	for _, kvp := range blockMap.OrderedKV {
		switch kvp.Key {
			case "timestamp":
				blockInfo.Timestamp, err = ojIntStrToUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("error parsing blockTimestamp: %w", err)
				}
			case "nonce":
				blockInfo.Nonce, err = ojIntStrToUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("error parsing blockNonce: %w", err)
				}
			case "round":
				blockInfo.Round, err = ojIntStrToUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("error parsing blockRound: %w", err)
				}
			case "epoch":
				blockInfo.Epoch, err = ojIntStrToUint64(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("error parsing blockEpoch: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown block info field: %s", kvp.Key)
		}
	}
	return blockInfo, nil
}

func (ae *Executor) setCurrentBlockInfo(blockInfo *BlockInfo) {
	ae.world.CurrentBlockInfo = &worldmock.BlockInfo{
		BlockTimestamp: blockInfo.Timestamp,
		BlockNonce:     blockInfo.Nonce,
		BlockRound:     blockInfo.Round,
		BlockEpoch:     uint32(blockInfo.Epoch),
		RandomSeed:     nil,
	}
}
