package main

import (
	"00097eth-relay/dao"
	"testing"
)

func TestBlockScanner_Start(t *testing.T) {
	nodeUrl := "https://ropsten.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	//nodeUrl := "https://mainnet.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	//nodeUrl := "https://kovan.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	//nodeUrl := "https://rinkeby.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	//nodeUrl := "https://goerli.infura.io/v3/00e4c7772d7f49399445dfaffc7b9a68"
	requester := NewETHRPCRequester(nodeUrl)
	option := dao.MySQLOptions{
		Hostname:           "127.0.0.1",
		Port:               "3306",
		User:               "root",
		Password:           "MySQL010801020",
		DBName:             "eth_relay",
		TablePrefix:        "eth_",
		MaxOpenConnections: 10,
		MaxIdleConnections: 5,
		ConnMaxLifetime:    15,
	}
	tables := []interface{}{dao.Block{}, dao.Transaction{}}
	mysqlConnector := dao.NewMySQLConnector(&option, tables)
	scanner := NewBlockScanner(*requester, mysqlConnector)
	err := scanner.Start()
	if err != nil {
		panic(err)
	}
	select {}
}
