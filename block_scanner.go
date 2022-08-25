package main

import (
	"00097eth-relay/dao"
	"00097eth-relay/model"
	"00097eth-relay/tool"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

type BlockScanner struct {
	ethRequester      ETHRPCRequester
	mysqlConnector    dao.MySQLConnector
	lastBlock         *dao.Block
	blockNumberToScan *big.Int
	fork              bool
	stop              chan bool
	lock              sync.Mutex
}

func NewBlockScanner(requester ETHRPCRequester, mysqlConnector dao.MySQLConnector) *BlockScanner {
	return &BlockScanner{
		ethRequester:   requester,
		mysqlConnector: mysqlConnector,
		lastBlock:      &dao.Block{},
		fork:           false,
		stop:           make(chan bool),
		lock:           sync.Mutex{},
	}
}
func (scanner *BlockScanner) forkCheck(currentBlock *dao.Block) bool {
	if currentBlock.BlockNumber == "" {
		panic("invalid block")
	}
	if scanner.lastBlock.BlockHash == currentBlock.BlockHash || scanner.lastBlock.BlockHash == currentBlock.ParentHash {
		scanner.lastBlock = currentBlock
		return false
	}
	startForkBlock, err := scanner.getStartForkBlock(currentBlock.ParentHash)
	if err != nil {
		panic(err)
	}
	scanner.lastBlock = startForkBlock
	numberEnd := scanner.hexToTen(currentBlock.BlockNumber).String()
	numberFrom := startForkBlock.BlockNumber
	_, err = scanner.mysqlConnector.DB.Table(dao.Block{}).Where("block_number > ? and block_number <= ?", numberFrom, numberEnd).Update(map[string]bool{"fork": true})
	if err != nil {
		panic(fmt.Errorf("update fork block failed, %s", err.Error()))
	}
	return true
}
func (scanner *BlockScanner) getStartForkBlock(parentHash string) (*dao.Block, error) {
	parent := dao.Block{}
	_, err := scanner.mysqlConnector.DB.Where("block_hash=?", parentHash).Get(&parent)
	if err == nil && parent.BlockNumber != "" {
		return &parent, nil
	}
	parentFull, err := scanner.retryGetBlockInfoByHash(parentHash)
	if err != nil {
		return nil, fmt.Errorf("severe fork error, restart block scan is needed. %s", err.Error())
	}
	return scanner.getStartForkBlock(parentFull.ParentHash)
}
func (scanner *BlockScanner) log(args ...interface{}) {
	fmt.Println(args...)
}
func (scanner *BlockScanner) retryGetBlockInfoByNumber(targetNumber *big.Int) (*model.FullBlock, error) {
Retry:
	fullBlock, err := scanner.ethRequester.GetBlockInfoByNumber(targetNumber)
	if err != nil {
		errInfo := err.Error()
		if strings.Contains(errInfo, "empty") {
			scanner.log("retrying...", targetNumber.String())
			goto Retry
		}
		return nil, err
	}
	return fullBlock, nil
}
func (scanner *BlockScanner) retryGetBlockInfoByHash(hash string) (*model.FullBlock, error) {
Retry:
	fullBlock, err := scanner.ethRequester.GetBlockInfoByHash(hash)
	if err != nil {
		errInfo := err.Error()
		if strings.Contains(errInfo, "empty") {
			scanner.log("retrying...", hash)
			goto Retry
		}
		return nil, err
	}
	return fullBlock, nil
}
func (scanner *BlockScanner) init() error {
	_, err := scanner.mysqlConnector.DB.Desc("create_time").Where("fork=?", false).Get(scanner.lastBlock)
	if err != nil {
		return err
	}
	if scanner.lastBlock.BlockHash == "" {
		latestBlockNumber, err := scanner.ethRequester.GetLatestBlockNumber()
		if err != nil {
			return err
		}
		latestBlock, err := scanner.ethRequester.GetBlockInfoByNumber(latestBlockNumber)
		if err != nil {
			return err
		}
		if latestBlock.Number == "" {
			panic(latestBlockNumber.String())
		}
		scanner.lastBlock.BlockHash = latestBlock.Hash
		scanner.lastBlock.ParentHash = latestBlock.ParentHash
		scanner.lastBlock.BlockNumber = latestBlock.Number
		scanner.lastBlock.CreateTime = scanner.hexToTen(latestBlock.Timestamp).Int64()
		scanner.blockNumberToScan = latestBlockNumber
	} else {
		scanner.blockNumberToScan, _ = new(big.Int).SetString(scanner.lastBlock.BlockNumber, 10)
		scanner.blockNumberToScan.Add(scanner.blockNumberToScan, new(big.Int).SetInt64(1))
	}
	return nil
}
func (scanner *BlockScanner) hexToTen(hex string) *big.Int {
	if strings.HasPrefix(hex, "0x") {
		ten, _ := new(big.Int).SetString(hex[2:], 16)
		return ten
	} else {
		ten, _ := new(big.Int).SetString(hex, 10)
		return ten
	}
}
func (scanner *BlockScanner) getScannerBlockNumber() (*big.Int, error) {
	latestBlockNumber, err := scanner.ethRequester.GetLatestBlockNumber()
	if err != nil {
		return nil, err
	}
	targetNumber := new(big.Int).Set(scanner.blockNumberToScan)
	if latestBlockNumber.Cmp(scanner.blockNumberToScan) < 0 {
	Next:
		for {
			select {
			case <-time.After(4 * time.Second):
				number, err := scanner.ethRequester.GetLatestBlockNumber()
				if err == nil && number.Cmp(scanner.blockNumberToScan) >= 0 {
					break Next
				}
			}
		}
	}
	return targetNumber, nil
}
func (scanner *BlockScanner) scan() error {
	targetNumber, err := scanner.getScannerBlockNumber()
	if err != nil {
		return err
	}
	fullBlock, err := scanner.retryGetBlockInfoByNumber(targetNumber)
	if err != nil {
		return err
	}
	scanner.blockNumberToScan.Add(scanner.blockNumberToScan, new(big.Int).SetInt64(1))
	tx := scanner.mysqlConnector.DB.NewSession()
	defer tx.Close()
	block := dao.Block{}
	_, err = tx.Where("block_hash=?", fullBlock.Hash).Get(&block)
	if err == nil && block.ID == 0 {
		block.BlockNumber = scanner.hexToTen(fullBlock.Number).String()
		block.ParentHash = fullBlock.ParentHash
		block.CreateTime = scanner.hexToTen(fullBlock.Timestamp).Int64()
		block.BlockHash = fullBlock.Hash
		block.Fork = false
		if _, err := tx.Insert(&block); err != nil {
			tool.ErrorPrintln(tx.Rollback(), "tx.Rollback() error")
			return err
		}
	}
	if scanner.forkCheck(&block) {
		data, _ := json.Marshal(fullBlock)
		scanner.log("fork occurred!", string(data))
		tool.ErrorPrintln(tx.Commit(), "tx.Commit() error")
		scanner.fork = true
		return errors.New("fork check")
	}
	scanner.log("scan block start ==> number:", scanner.hexToTen(fullBlock.Number), "hash:", fullBlock.Hash)
	for index, transaction := range fullBlock.Transactions {
		scanner.log("tx hash ==> ", transaction.Hash)
		if index == 5 {
			break
		}
	}
	scanner.log("scan block finished\n===========================")
	if _, err = tx.Insert(&fullBlock.Transactions); err != nil {
		tool.ErrorPrintln(tx.Rollback(), "tx.Rollback() error")
		return err
	}
	return tx.Commit()
}
func (scanner *BlockScanner) Start() error {
	scanner.lock.Lock()
	if err := scanner.init(); err != nil {
		scanner.lock.Unlock()
		return err
	}
	executeScan := func() {
		if err := scanner.scan(); err != nil {
			scanner.log(err.Error())
			return
		}
		time.Sleep(1 * time.Second)
	}
	go func() {
		for {
			select {
			case <-scanner.stop:
				scanner.log("finish block scanner!")
				return
			default:
				if !scanner.fork {
					executeScan()
				} else {
					if err := scanner.init(); err != nil {
						scanner.log(err.Error())
						return
					}
					scanner.fork = false
				}
			}
		}
	}()
	return nil
}
func (scanner *BlockScanner) Stop() {
	scanner.lock.Unlock()
	scanner.stop <- true
}
