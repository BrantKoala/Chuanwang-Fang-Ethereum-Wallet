package tool

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
)

func TestMakeMethodId(t *testing.T) {
	contractABI := `[
	{
		"inputs": [
			{
				"internalType": "uint8",
				"name": "arg1",
				"type": "uint8"
			},
			{
				"internalType": "uint8",
				"name": "arg2",
				"type": "uint8"
			}
		],
		"name": "fangChuanWangCalculate",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "pure",
		"type": "function"
	}
]`
	methodName := "fangChuanWangCalculate"
	methodId, err := MakeMethodId(methodName, contractABI)
	if err != nil {
		fmt.Println("generate methodId failed, ", err.Error())
		return
	}
	fmt.Println("methodId is ", methodId)
}

func TestUnlockETHWallet(t *testing.T) {
	address := "0x3d872c4dc3EAF7EEA62d5C3762fD6d4BB524B0E7"
	keysDir := "../keystore"
	err1 := UnlockETHWallet(keysDir, address, "789")
	if err1 != nil {
		fmt.Println("No1. unlock error", err1.Error())
	} else {
		fmt.Println("No1. unlock success")
	}
	err2 := UnlockETHWallet(keysDir, address, "123456abc")
	if err2 != nil {
		fmt.Println("No2. unlock error", err2.Error())
	} else {
		fmt.Println("No2. unlock success")
	}
}
func TestSignETHTransaction(t *testing.T) {
	address := "0x3d872c4dc3EAF7EEA62d5C3762fD6d4BB524B0E7"
	keysDir := "../keystore"
	err1 := UnlockETHWallet(keysDir, address, "789")
	if err1 != nil {
		fmt.Println("No1. unlock error", err1.Error())
	} else {
		fmt.Println("No1. unlock success")
	}
	err2 := UnlockETHWallet(keysDir, address, "123456abc")
	if err2 != nil {
		fmt.Println("No2. unlock error", err2.Error())
	} else {
		fmt.Println("No2. unlock success")
	}
	tx := types.NewTransaction(
		123,
		common.Address{},
		new(big.Int).SetInt64(10),
		1000,
		new(big.Int).SetInt64(20),
		[]byte("transaction"))
	signTx, err := SignETHTransaction(address, tx)
	if err != nil {
		fmt.Println("sign failed", err.Error())
		return
	}
	data, _ := json.Marshal(signTx)
	fmt.Println("sign success", string(data))
}
