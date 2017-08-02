package main

import (
	"fmt"
	"scheduler/db"
	"scheduler/model"
	"scheduler/webService"
	"time"
)

func main() {
	para := &db.DbConnPara{
		Host:     "192.168.56.132",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "request"}

	cephManager := model.NewCephManager(para)
	if nil == cephManager {
		fmt.Println("Start SMP fail!")
		return
	}

	webService.SetupWebService(cephManager)

	for {
		time.Sleep(600 * time.Second)
	}
}
