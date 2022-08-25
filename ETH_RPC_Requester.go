package main

import (
	"00097eth-relay/model"
	"00097eth-relay/tool"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type ETHRPCRequester struct {
	nonceManager *NonceManager
	client       *ETHRPCClient
}

func NewETHRPCRequester(nodeUrl string) *ETHRPCRequester {
	return &ETHRPCRequester{nonceManager: NewNonceManager(), client: NewETHRPCClient(nodeUrl)}
}
func (r *ETHRPCRequester) GetTransactionByHash(txHash string) (model.Transaction, error) {
	methodName := "eth_getTransactionByHash"
	result := model.Transaction{}
	err := r.client.GetRpc().Call(&result, methodName, txHash)
	return result, err
}
func (r *ETHRPCRequester) GetTransactions(txHashs []string) ([]*model.Transaction, error) {
	name := "eth_getTransactionByHash"
	rets := []*model.Transaction{}
	size := len(txHashs)
	reqs := []rpc.BatchElem{}
	for i := 0; i < size; i++ {
		ret := model.Transaction{}
		req := rpc.BatchElem{
			Method: name,
			Args:   []interface{}{txHashs[i]},
			Result: &ret,
		}
		reqs = append(reqs, req)
		rets = append(rets, &ret)
	}
	err := r.client.GetRpc().BatchCall(reqs)
	return rets, err
}
func (r *ETHRPCRequester) GetETHBalance(address string) (string, error) {
	name := "eth_getBalance"
	result := ""
	err := r.client.GetRpc().Call(&result, name, address, "latest")
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("eth balance is null")
	}
	ten, _ := new(big.Int).SetString(result[2:], 16)
	return ten.String(), err
}
func (r *ETHRPCRequester) GetETHBalances(addresses []string) ([]string, error) {
	name := "eth_getBalance"
	rets := []*string{}
	size := len(addresses)
	reqs := []rpc.BatchElem{}
	for i := 0; i < size; i++ {
		ret := ""
		req := rpc.BatchElem{
			Method: name,
			Args:   []interface{}{addresses[i], "latest"},
			Result: &ret,
		}
		reqs = append(reqs, req)
		rets = append(rets, &ret)
	}
	err := r.client.GetRpc().BatchCall(reqs)
	if err != nil {
		return nil, err
	}
	for _, req := range reqs {
		if req.Error != nil {
			return nil, req.Error
		}
	}
	finalRet := []string{}
	for _, item := range rets {
		ten, _ := new(big.Int).SetString((*item)[2:], 16)
		finalRet = append(finalRet, ten.String())
	}
	return finalRet, err
}

type ERC20BalanceRpcReq struct {
	ContractAddress string
	UserAddress     string
	ContractDecimal int
}

func (r *ETHRPCRequester) GetERC20Balances(paramArr []ERC20BalanceRpcReq) ([]string, error) {
	name := "eth_call"
	methodId := "0x70a08231"
	rets := []*string{}
	size := len(paramArr)
	reqs := []rpc.BatchElem{}

	for i := 0; i < size; i++ {
		ret := ""
		arg := &model.CallArg{}
		userAddress := paramArr[i].UserAddress
		arg.Gas = hexutil.EncodeUint64(300000)
		arg.To = common.HexToAddress(paramArr[i].ContractAddress)
		arg.Data = methodId + "000000000000000000000000" + userAddress[2:]
		req := rpc.BatchElem{
			Method: name,
			Args:   []interface{}{arg, "latest"},
			Result: &ret,
		}
		reqs = append(reqs, req)
		rets = append(rets, &ret)
	}
	err := r.client.GetRpc().BatchCall(reqs)
	if err != nil {
		return nil, err
	}
	for _, req := range reqs {
		if req.Error != nil {
			return nil, req.Error
		}
	}
	finalRet := []string{}
	for _, item := range rets {
		if *item == "" {
			//TODO: This probably needs to change
			continue
		}
		ten, _ := new(big.Int).SetString((*item)[2:], 16)
		finalRet = append(finalRet, ten.String())
	}
	return finalRet, err
}
func (r *ETHRPCRequester) GetLatestBlockNumber() (*big.Int, error) {
	methodName := "eth_blockNumber"
	number := ""
	err := r.client.GetRpc().Call(&number, methodName)
	if err != nil {
		return nil, fmt.Errorf("get latest block number failed! %s", err.Error())
	}
	ten, _ := new(big.Int).SetString(number[2:], 16)
	return ten, nil
}

func (r *ETHRPCRequester) GetBlockInfoByNumber(blockNumber *big.Int) (*model.FullBlock, error) {
	number := fmt.Sprintf("%#x", blockNumber)
	methodName := "eth_getBlockByNumber"
	fullBlock := model.FullBlock{}
	err := r.client.GetRpc().Call(&fullBlock, methodName, number, true)
	if err != nil {
		return nil, fmt.Errorf("get block info failed! %s", err.Error())
	}
	if fullBlock.Number == "" {
		return nil, fmt.Errorf("block info is empty, %s", err.Error())
	}
	return &fullBlock, nil
}

func (r *ETHRPCRequester) GetBlockInfoByHash(blockHash string) (*model.FullBlock, error) {
	methodName := "eth_getBlockByHash"
	fullBlock := model.FullBlock{}
	err := r.client.GetRpc().Call(&fullBlock, methodName, blockHash, true)
	if err != nil {
		return nil, fmt.Errorf("get block info failed, %s", err.Error())
	}
	if fullBlock.Number == "" {
		return nil, fmt.Errorf("block info is empty %s", blockHash)
	}
	return &fullBlock, nil
}
func (r *ETHRPCRequester) ETHCall(result interface{}, arg model.CallArg) error {
	methodName := "eth_call"
	err := r.client.GetRpc().Call(result, methodName, arg, "latest")
	if err != nil {
		return fmt.Errorf("eth_call failed, %s", err.Error())
	}
	return nil
}
func (r *ETHRPCRequester) CreateETHWallet(password string) (string, error) {
	if password == "" {
		return "", errors.New("password can't be empty")
	}
	if len(password) < 6 {
		return "", errors.New("password's length must be no less than 6")
	}
	keyDir := "./keystore"
	ks := keystore.NewKeyStore(keyDir, keystore.StandardScryptN, keystore.StandardScryptP)
	wallet, err := ks.NewAccount(password)
	if err != nil {
		return "0x", err
	}
	return wallet.Address.String(), err
}
func (r *ETHRPCRequester) SendTransaction(address string, transaction *types.Transaction) (string, error) {
	signedTx, err := tool.SignETHTransaction(address, transaction)
	if err != nil {
		return "", fmt.Errorf("sign failed, %s", err.Error())
	}
	txRlpData, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", fmt.Errorf("rlp serialization error, %s", err.Error())
	}
	txHash := ""
	methodName := "eth_sendRawTransaction"
	//TODO: hexutil.Encode maybe incorrect
	err = r.client.GetRpc().Call(&txHash, methodName, hexutil.Encode(txRlpData))
	if err != nil {
		return "", fmt.Errorf("send transaction failed, %s", err.Error())
	}
	oldNonce := r.nonceManager.GetNonce(address)
	if oldNonce == nil {
		r.nonceManager.SetNonce(address, new(big.Int).SetUint64(transaction.Nonce()))
	}
	r.nonceManager.PlusNonce(address)
	return txHash, nil
}
func (r *ETHRPCRequester) GetNonce(address string) (uint64, error) {
	methodName := "eth_getTransactionCount"
	nonce := ""
	err := r.client.GetRpc().Call(&nonce, methodName, address, "pending")
	if err != nil {
		return 0, fmt.Errorf("get nonce failed, %s", err.Error())
	}
	n, _ := new(big.Int).SetString(nonce[2:], 16)
	return n.Uint64(), nil
}
func (r *ETHRPCRequester) SendETHTransaction(fromStr, toStr, valueStr string, gasLimit, gasPrice uint64) (string, error) {
	if !common.IsHexAddress(fromStr) || !common.IsHexAddress(toStr) {
		return "", errors.New("invalid address")
	}
	to := common.HexToAddress(toStr)
	TxGasPrice := new(big.Int).SetUint64(gasPrice)
	realValue := tool.GetRealDecimalValue(valueStr, 18)
	if realValue == "" {
		return "", errors.New("invalid value")
	}
	amount, _ := new(big.Int).SetString(realValue, 10)
	nonce := r.nonceManager.GetNonce(fromStr)
	if nonce == nil {
		n, err := r.GetNonce(fromStr)
		if err != nil {
			return "", fmt.Errorf("get nonce failed %s", err.Error())
		}
		nonce = new(big.Int).SetUint64(n)
		r.nonceManager.SetNonce(fromStr, nonce)
	}
	data := []byte("")
	transaction := types.NewTransaction(
		nonce.Uint64(),
		to,
		amount,
		gasLimit,
		TxGasPrice,
		data)
	return r.SendTransaction(fromStr, transaction)
}
func (r *ETHRPCRequester) SendERC20Transaction(fromStr, contract, receiver, valueStr string, gasLimit, gasPrice uint64, decimal int) (string, error) {
	if !common.IsHexAddress(fromStr) || !common.IsHexAddress(contract) || !common.IsHexAddress(receiver) {
		return "", errors.New("invalid address")
	}
	to := common.HexToAddress(contract)
	txGasPrice := new(big.Int).SetUint64(gasPrice)
	amount := new(big.Int).SetInt64(0)
	nonce := r.nonceManager.GetNonce(fromStr)
	if nonce == nil {
		n, err := r.GetNonce(fromStr)
		if err != nil {
			return "", fmt.Errorf("get nonce failed %s", err.Error())
		}
		nonce = new(big.Int).SetUint64(n)
		r.nonceManager.SetNonce(fromStr, nonce)
	}
	data := tool.BuildERC20TransferData(receiver, valueStr, decimal)
	dataBytes := common.FromHex(data)
	transaction := types.NewTransaction(
		nonce.Uint64(),
		to,
		amount,
		gasLimit,
		txGasPrice,
		dataBytes)
	return r.SendTransaction(fromStr, transaction)
}
