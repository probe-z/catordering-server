package main

import (
	"net/http"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// server configuration
var conf struct {
	Server struct {
		Addr                string `json:"addr"`
		WriteTimeout        int    `json:"writeTimeout"`
		ReadTimeout         int    `json:"readTimeout"`
		IdleTimeout         int    `json:"idleTimeout"`
		RestartWaitDuration int    `json:"restartWaitDuration"`
	} `json:"server"`
	Database struct {
		DBConf
		Database   string `json:"database"`
		UserTable  string `json:"userTable"`
		FoodTable  string `json:"foodTable"`
		OrderTable string `json:"orderTable"`
	} `json:"database"`
}

var db *sql.DB

func main() {
	err := ParseJsonConf(&conf, "conf.json")
	if err != nil {
		return
	}

	db, err = NewDB(&conf.Database.DBConf)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	handler := mux.NewRouter()
	handler.HandleFunc("/food", getFoodListHandler).Methods("GET")
	srv := &http.Server{
		Addr:         conf.Server.Addr,
		WriteTimeout: time.Second * time.Duration(conf.Server.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(conf.Server.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(conf.Server.IdleTimeout),
		Handler:      handler,
	}
	StartServer(srv, conf.Server.RestartWaitDuration)
}
