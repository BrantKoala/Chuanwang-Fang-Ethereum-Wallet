package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
	"xorm.io/core"
)

type MySQLOptions struct {
	Hostname           string
	Port               string
	User               string
	Password           string
	DBName             string
	TablePrefix        string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    int
}
type MySQLConnector struct {
	options *MySQLOptions
	tables  []interface{}
	DB      *xorm.Engine
}

func (s *MySQLConnector) createTables() error {
	if len(s.tables) == 0 {
		return nil
	}
	if err := s.DB.CreateTables(s.tables...); err != nil {
		return fmt.Errorf("create mysql table error: %s", err.Error())
	}
	if err := s.DB.Sync2(s.tables...); err != nil {
		return fmt.Errorf("sync table error: %s", err.Error())
	}
	return nil
}
func NewMySQLConnector(options *MySQLOptions, tables []interface{}) MySQLConnector {
	var connector MySQLConnector
	connector.options = options
	connector.tables = tables
	url := ""
	if options.Hostname == "" || options.Hostname == "127.0.0.1" {
		url = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", options.User, options.Password, options.DBName)
	} else {
		url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", options.User, options.Password, options.Hostname, options.Port, options.DBName)
	}
	db, err := xorm.NewEngine("mysql", url)
	if err != nil {
		panic(fmt.Errorf("database initialization failed, %s", err.Error()))
	}
	tableMapper := core.NewPrefixMapper(core.SnakeMapper{}, options.TablePrefix)
	db.SetTableMapper(tableMapper)
	db.DB().SetConnMaxLifetime(time.Duration(options.ConnMaxLifetime) * time.Second)
	db.DB().SetMaxIdleConns(options.MaxIdleConnections)
	db.DB().SetMaxOpenConns(options.MaxOpenConnections)
	//db.ShowSQL(true)
	if err = db.Ping(); err != nil {
		panic(fmt.Errorf("database connection failed, %s", err.Error()))
	}
	connector.DB = db
	if err = connector.createTables(); err != nil {
		panic(fmt.Errorf("create tables failed, %s", err.Error()))
	}
	return connector
}
