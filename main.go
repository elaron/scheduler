package main

import (
	"fmt"
	"scheduler/auth"
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
		fmt.Println("Start SMP CephManager fail!")
		return
	}

	para_auth := &db.DbConnPara{
		Host:     "192.168.56.132",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "auth"}

	a := auth.NewAuthManager(para_auth)
	if nil == cephManager {
		fmt.Println("Start SMP AuthManager fail!")
		return
	}

	webService.SetupWebService(cephManager, a)

	for {
		time.Sleep(600 * time.Second)
	}
}
