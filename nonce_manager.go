package main

import (
	"math/big"
	"sync"
)

type NonceManager struct {
	lock          sync.Mutex
	nonceMemCache map[string]*big.Int
}

func NewNonceManager() *NonceManager {
	return &NonceManager{lock: sync.Mutex{}}
}
func (n *NonceManager) SetNonce(address string, nonce *big.Int) {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	n.nonceMemCache[address] = nonce
}
func (n *NonceManager) GetNonce(address string) *big.Int {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.nonceMemCache[address]
}
func (n *NonceManager) PlusNonce(address string) {
	if n.nonceMemCache == nil {
		n.nonceMemCache = map[string]*big.Int{}
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	oldNonce := n.nonceMemCache[address]
	newNonce := oldNonce.Add(oldNonce, big.NewInt(int64(1)))
	n.nonceMemCache[address] = newNonce
}
