package scdb

import (
	"database/sql"
	"log"

	"github.com/snxamdf/scgo/data"
)

type DBSourceInterface interface {
	DB() *sql.DB
	DataBaseType() string
}

type Config struct {
	DriverName                 string
	Alias, UserName, PassWord  string
	Ip, Prot, DBName, Charset  string
	MaxIdleConns, MaxOpenConns int
	Db                         *sql.DB
	dataBaseType               string
}

func (this *Config) DB() *sql.DB {
	return this.Db
}

func (this *Config) DataBaseType() string {
	return this.dataBaseType
}

func (this *Config) MySqlInit() error {
	if this.Charset == "" {
		this.Charset = "UTF8"
	}
	var dataSource = this.UserName + ":" + this.PassWord + "@tcp(" + this.Ip + ":" + this.Prot + ")/" + this.DBName + "?charset=" + this.Charset
	log.Println("data source [", dataSource, "]")
	db, err := sql.Open(this.DriverName, dataSource)
	if err != nil {
		log.Println(err)
		return err
	}

	db.SetMaxIdleConns(this.MaxIdleConns)
	db.SetMaxOpenConns(this.MaxOpenConns)
	this.Db = db
	this.dataBaseType = data.DATA_BASE_MYSQL
	return nil
}

func (this *Config) OracleInit() error {
	if this.Charset == "" {
		this.Charset = "UTF8"
	}
	var dataSource = this.UserName + ":" + this.PassWord + "@tcp(" + this.Ip + ":" + this.Prot + ")/" + this.DBName + "?charset=" + this.Charset
	log.Println("data source :", dataSource)
	db, err := sql.Open(this.DriverName, dataSource)
	if err != nil {
		log.Println(err)
		return err
	}

	db.SetMaxIdleConns(this.MaxIdleConns)
	db.SetMaxOpenConns(this.MaxOpenConns)
	this.Db = db
	this.dataBaseType = data.DATA_BASE_ORACLE
	return nil
}
