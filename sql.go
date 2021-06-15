package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBConf struct {
	DSN          string `json:"dsn"`
	MaxOpenConns int    `json:"maxOpenConns"`
	MaxIdleConns int    `json:"maxIdleConns"`
}

func NewDB(dbConf *DBConf) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbConf.DSN)
	if err != nil {
		log.Printf("sql open error:%v\n", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("sql ping error:%v\n", err)
		return nil, err
	}
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	return db, nil
}
