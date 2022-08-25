package dao

import (
	"fmt"
	"testing"
)

func TestNewMySQLConnector(t *testing.T) {
	option := MySQLOptions{
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
	tables := []interface{}{}
	tables = append(tables, Block{}, Transaction{})
	mysql := NewMySQLConnector(&option, tables)
	if mysql.DB.Ping() == nil {
		fmt.Println("database connection success")
		fmt.Println("create table success")
	} else {
		fmt.Println("database connection failed")
	}
}
