package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"
)

func bigintToStr(value *big.Int) string {
	return fmt.Sprintf("%d", value)
}

func intToStr(value int) string {
	return fmt.Sprintf("%d", value);
}

func strToOjStr(str string) oj.OJsonObject {
	return &oj.OJsonString{Value: str}
}

func bigintToOjIntStr(value *big.Int) oj.OJsonObject {
	return strToOjStr(bigintToStr(value))
}

func intToOjIntStr(value int) oj.OJsonObject {
	return strToOjStr(intToStr(value))
}

func bytesToOjHexStr(bytes []byte) oj.OJsonObject {
	return strToOjStr(hex.EncodeToString(bytes))
}

func ojStrToStr(obj oj.OJsonObject) (string, error) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return "", errors.New("not a string value")
	}
	return str.Value, nil
}

func ojBoolToBool(obj oj.OJsonObject) (bool, error) {
	value, isBool := obj.(*oj.OJsonBool)
	if !isBool {
		return false, errors.New("not a bool value")
	}
	return bool(*value), nil
}

func ojStrToBytes(obj oj.OJsonObject) ([]byte, error) {
	str, err := ojStrToStr(obj)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func ojHexStrToBytes(obj oj.OJsonObject) ([]byte, error) {
	strVal, err := ojStrToStr(obj)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strVal)
}

func ojIntStrToBigint(obj oj.OJsonObject) (*big.Int, error) {
	strVal, err := ojStrToStr(obj)
	if err != nil {
		return nil, err
	}
	n, isInt := big.NewInt(0).SetString(strVal, 10)
	if !isInt {
		return nil, errors.New("not a bigint")
	}
	return n, nil
}

func ojIntStrToUint64(obj oj.OJsonObject) (uint64, error) {
	bi, err := ojIntStrToBigint(obj)
	if err != nil {
		return 0, err
	}
	if bi == nil || !bi.IsUint64() {
		return 0, errors.New("value is not uint64")
	}
	return bi.Uint64(), nil
}
