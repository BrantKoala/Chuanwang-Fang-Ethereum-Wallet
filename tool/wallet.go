package tool

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func MakeMethodId(methodName string, abiStr string) (string, error) {
	abiObj := &abi.ABI{}
	err := abiObj.UnmarshalJSON([]byte(abiStr))
	if err != nil {
		return "", err
	}
	method := abiObj.Methods[methodName]
	methodIdBytes := method.ID
	methodId := "0x" + common.Bytes2Hex(methodIdBytes)
	return methodId, nil
}

var ETHUnlockMap map[string]accounts.Account
var UnlockKeystore *keystore.KeyStore

func UnlockETHWallet(keysDir string, address, password string) error {
	if UnlockKeystore == nil {
		UnlockKeystore = keystore.NewKeyStore(keysDir, keystore.StandardScryptN, keystore.StandardScryptP)
		if UnlockKeystore == nil {
			return errors.New("keystore is nil")
		}
	}
	accountToUnlock := accounts.Account{Address: common.HexToAddress(address)}
	if err := UnlockKeystore.Unlock(accountToUnlock, password); nil != err {
		return errors.New("unlock err: " + err.Error())
	}
	if ETHUnlockMap == nil {
		ETHUnlockMap = map[string]accounts.Account{}
	}
	ETHUnlockMap[address] = accountToUnlock
	return nil
}

type txData struct {
	AccountNonce uint64          `json:"nonce" gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas" gencodec:"required"`
	Recipient    *common.Address `json:"to" rlp:"nil"`
	Amount       *big.Int        `json:"value" gencodec:"required"`
	Payload      []byte          `json:"input" gencodec:"required"`
	V            *big.Int        `json:"v" gencodec:"required"`
	R            *big.Int        `json:"r" gencodec:"required"`
	S            *big.Int        `json:"s" gencodec:"required"`
	Hash         *common.Hash    `json:"hash" rlp:"-"`
}

func SignETHTransaction(address string, transaction *types.Transaction) (*types.Transaction, error) {
	if UnlockKeystore == nil {
		return nil, errors.New("initialize keystore is needed")
	}
	account := ETHUnlockMap[address]
	if !common.IsHexAddress(account.Address.String()) {
		return nil, errors.New("account need to unlock first")
	}
	return UnlockKeystore.SignTx(account, transaction, nil)
}
