package main

import (
	"context"
	"flag"
	"fmt"
	"scheduler/auth"
	"scheduler/db"
	"scheduler/model"
	"scheduler/webService"
	"time"
)

type SysConfig struct {
	cephManagerDbIp   string
	cephManagerDbPort int32

	authManagerDbIp   string
	authManagerDbPort int32
}

func initConfig(config *SysConfig) {
	ip1 := flag.String("cephMngDbIp", "0.0.0.0", "Ceph manager db ip address")
	port1 := flag.Int("cephMngDbPort", 5432, "Ceph manager db port")

	ip2 := flag.String("authMngDbIp", "0.0.0.0", "Auth manager db ip address")
	port2 := flag.Int("authMngDbPort", 5432, "Auth manager db port")

	flag.Parse()

	config.cephManagerDbIp = *ip1
	config.cephManagerDbPort = int32(*port1)
	config.authManagerDbIp = *ip2
	config.authManagerDbPort = int32(*port2)
}

func main() {

	var conf SysConfig
	initConfig(&conf)

	para_request := &db.DbConnPara{
		Host:     conf.cephManagerDbIp,
		Port:     conf.cephManagerDbPort,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "request"}

	cephManager := model.NewCephManager(para_request)
	if nil == cephManager {
		fmt.Println("Start SMP CephManager fail!")
		return
	}

	para_auth := &db.DbConnPara{
		Host:     conf.authManagerDbIp,
		Port:     conf.authManagerDbPort,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "auth"}

	a := auth.NewAuthManager(para_auth)
	if nil == cephManager {
		fmt.Println("Start SMP AuthManager fail!")
		return
	}

	//ctx, cancel := context.WithCancel(context.Background())
	ctx, _ := context.WithCancel(context.Background())

	webService.SetupWebService(ctx, cephManager, a)

	for {
		time.Sleep(600 * time.Second)
	}
}
