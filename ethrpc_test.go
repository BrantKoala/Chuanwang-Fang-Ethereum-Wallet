package main

import (
	"00097eth-relay/model"
	"00097eth-relay/tool"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"testing"
)

func TestNewETHRPClient(t *testing.T) {
	test2 := NewETHRPCClient("www.nihao.com").GetRpc()
	if test2 == nil {
		fmt.Println("Initialization failed.")
	}
	client := NewETHRPCClient("123://456").GetRpc()
	if client == nil {
		fmt.Println("Initialization failed.")
	}
}
func TestETHRPCRequester_GetTransactionByHash(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	txHash := "0xd8380c84227179928b1750d3e02fddb4b2696b0420b28606f2b703c482a071ba"
	if txHash == "" || len(txHash) != 66 {
		fmt.Println("Illegal transaction hash")
		return
	}
	txInfo, err := NewETHRPCRequester(nodeUrl).GetTransactionByHash(txHash)
	if err != nil {
		fmt.Println("query transaction failed, information: ", err.Error())
		return
	}
	jsonData, _ := json.Marshal(txInfo)
	fmt.Println(string(jsonData))
}
func TestETHRPCRequester_GetTransactions(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	txHash1 := "0x90d877f84f5771dfaa8145aee15f3077405db45035aebe431da5d68f277f945c"
	txHash2 := "0x90d877f84f5771dfaa8145aee15f3077405db45035aebe431da5d68f277f945b"
	txHash3 := "0x619436ebae899ec9208b1fede95e3845d94e1441eaf3d07bc42c12a19747c437"
	txHashes := []string{txHash1, txHash2, txHash3}
	if txHashes == nil || len(txHashes) == 0 {
		fmt.Println("Illegal hash array")
		return
	}
	txInfos, err := NewETHRPCRequester(nodeUrl).GetTransactions(txHashes)
	if err != nil {
		fmt.Println("query transaction failed, info: ", err.Error())
		return
	}
	jsonData, _ := json.Marshal(txInfos)
	fmt.Println(string(jsonData))
}
func TestETHRPCRequester_GetETHBalance(t *testing.T) {
	//nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	//nodeUrl := "https://mainnet.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	nodeUrl := "https://kovan.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	address := "0xcC4DbB4208142EB1A061F87bf4cB71E1D0Aa13F2"
	if address == "" || len(address) != 42 {
		fmt.Println("Illegal address")
		return
	}
	balance, err := NewETHRPCRequester(nodeUrl).GetETHBalance(address)
	if err != nil {
		fmt.Println("query eth balance failed, info: ", err.Error())
		return
	}
	fmt.Println(balance)
}

func TestETHRPCRequester_GetETHBalances(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	addresses := []string{
		"0x56bc0bb5359BfDF9BADc0b4F1E76669Ab88D17a7",
		"0xb12713bfa9d1de339ca14b01f8f14f092ffe75bf",
	}
	balance, err := NewETHRPCRequester(nodeUrl).GetETHBalances(addresses)
	if err != nil {
		fmt.Println("query eth failed, info: ", err.Error())
		return
	}
	fmt.Println(balance)
}
func TestETHRPCRequester_GetERC20Balances(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	paramsGEB := []ERC20BalanceRpcReq{
		{
			ContractAddress: "0x29365197AFD915d49293699C49796122A7EDba8c",
			UserAddress:     "0xcC4DbB4208142EB1A061F87bf4cB71E1D0Aa13F2",
			ContractDecimal: 18,
		}, {
			ContractAddress: "0x29365197AFD915d49293699C49796122A7EDba8c",
			UserAddress:     "0x3d872c4dc3EAF7EEA62d5C3762fD6d4BB524B0E7",
			ContractDecimal: 18,
		},
	}
	balance, err := NewETHRPCRequester(nodeUrl).GetERC20Balances(paramsGEB)
	if err != nil {
		fmt.Println("query ERC20 token failed, info: ", err.Error())
		return
	}
	fmt.Println(balance)
}
func TestETHRPCRequester_GetLatestBlockNumber(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	number, err := NewETHRPCRequester(nodeUrl).GetLatestBlockNumber()
	if err != nil {
		fmt.Println("get latest block number failed, info: ", err.Error())
		return
	}
	fmt.Println("decimalism: ", number.String())
}

func TestETHRPCRequester_GetBlockInfoByNumber(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	requester := NewETHRPCRequester(nodeUrl)
	number, _ := requester.GetLatestBlockNumber()
	fmt.Println("Block number is ", number)
	fullBlock, err := requester.GetBlockInfoByNumber(number)
	if err != nil {
		fmt.Println("get block info failed, ", err.Error())
		return
	}
	jsonData, _ := json.Marshal(fullBlock)
	fmt.Println("get block info by number", string(jsonData))
}
func TestETHRPCRequester_GetBlockInfoByHash(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	requester := NewETHRPCRequester(nodeUrl)
	blockHash := "0x3e21d56afe7c30f1434cb8ab271fad7fca22844a4f97fcdfea5eea4a9478619c"
	fullBlock, err := requester.GetBlockInfoByHash(blockHash)
	if err != nil {
		fmt.Println("get block info failed, ", err.Error())
		return
	}
	jsonData, _ := json.Marshal(fullBlock)
	fmt.Println("get block info by hash: ", string(jsonData))
}
func TestETHRPCRequester_ETHCall(t *testing.T) {
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
	methodId, err := tool.MakeMethodId(methodName, contractABI)
	if err != nil {
		panic(err)
	}
	arg1 := common.HexToHash("3").String()[2:]
	arg2 := common.HexToHash("4").String()[2:]
	contractAddress := "0x575622edB6409a64c4f69cf1df82F2689C09c920"
	args := model.CallArg{
		To:   common.HexToAddress(contractAddress),
		Gas:  hexutil.EncodeUint64(300000),
		Data: methodId + arg1 + arg2,
	}
	result := ""
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	err = NewETHRPCRequester(nodeUrl).ETHCall(&result, args)
	if err != nil {
		panic(err)
	}
	ten, _ := new(big.Int).SetString(result[2:], 16)
	fmt.Println("the result is ", ten.String())
}
func TestETHRPCRequester_CreateETHWallet(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	address1, err := NewETHRPCRequester(nodeUrl).CreateETHWallet("12345")
	if err != nil {
		fmt.Println("No.1 create failed, ", err.Error())
	} else {
		fmt.Println("No.1 create success, ", address1)
	}
	address2, err := NewETHRPCRequester(nodeUrl).CreateETHWallet("123456abc")
	if err != nil {
		fmt.Println("No.2 create failed, ", err.Error())
	} else {
		fmt.Println("No.2 create success, ", address2)
	}
}
func TestETHRPCRequester_GetNonce(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	address := "0xcC4DbB4208142EB1A061F87bf4cB71E1D0Aa13F2"
	if address == "" || len(address) != 42 {
		fmt.Println("Illegal transaction address")
		return
	}
	nonce, err := NewETHRPCRequester(nodeUrl).GetNonce(address)
	if err != nil {
		fmt.Println("query nonce failed, ", err.Error())
		return
	}
	fmt.Println(nonce)
}
func TestETHRPCRequester_SendETHTransaction(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	from := "0x3d872c4dc3EAF7EEA62d5C3762fD6d4BB524B0E7"
	if from == "" || len(from) != 42 {
		fmt.Println("Illegal address")
		return
	}
	to := "0xcC4DbB4208142EB1A061F87bf4cB71E1D0Aa13F2"
	value := "0.1"
	gasLimit := uint64(100000)
	gasPrice := uint64(36000000000)
	err := tool.UnlockETHWallet("./keystore", from, "123456abc")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	txHash, err := NewETHRPCRequester(nodeUrl).SendETHTransaction(from, to, value, gasLimit, gasPrice)
	if err != nil {
		fmt.Println("ETH transaction failed, ", err.Error())
		return
	}
	fmt.Println(txHash)
}

func TestETHRPCRequester_SendERC20Transaction(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	from := "0x3d872c4dc3EAF7EEA62d5C3762fD6d4BB524B0E7"
	if from == "" || len(from) != 42 {
		fmt.Println("Illegal transaction address")
		return
	}
	to := "0x29365197AFD915d49293699C49796122A7EDba8c"
	amount := "10"
	decimal := 18
	receiver := "0xcC4DbB4208142EB1A061F87bf4cB71E1D0Aa13F2"
	gasLimit := uint64(100000)
	gasPrice := uint64(24000000000)
	err := tool.UnlockETHWallet("./keystore", from, "123456abc")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	txHash, err := NewETHRPCRequester(nodeUrl).SendERC20Transaction(from, to, receiver, amount, gasLimit, gasPrice, decimal)
	if err != nil {
		fmt.Println("ERC20 token transfer failed, ", err.Error())
		return
	}
	fmt.Println(txHash)
}
