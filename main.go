package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"scheduler/auth"
	"scheduler/db"
	"scheduler/model"
	"scheduler/webService"
	"sync"
)

type SysConfig struct {
	cephManagerDbIp   string
	cephManagerDbPort int32

	authManagerDbIp   string
	authManagerDbPort int32

	elasticsearchManagerDbIp   string
	elasticsearchManagerDbPort int32

	logDir string
}

func initConfig(config *SysConfig) {
	ip1 := flag.String("cephMngDbIp", "0.0.0.0", "Ceph manager db ip address")
	port1 := flag.Int("cephMngDbPort", 5432, "Ceph manager db port")

	ip2 := flag.String("authMngDbIp", "0.0.0.0", "Auth manager db ip address")
	port2 := flag.Int("authMngDbPort", 5432, "Auth manager db port")

	ip3 := flag.String("esMngDbIp", "0.0.0.0", "Elasticsearch manager db ip address")
	port3 := flag.Int("esMngDbPort", 9200, "Elasticsearch manager db port")

	logDir := flag.String("logPath", "/var/log/sheduler", "SMP log path")

	flag.Parse()

	config.cephManagerDbIp = *ip1
	config.cephManagerDbPort = int32(*port1)
	config.authManagerDbIp = *ip2
	config.authManagerDbPort = int32(*port2)
	config.elasticsearchManagerDbIp = *ip3
	config.elasticsearchManagerDbPort = int32(*port3)
	config.logDir = *logDir
}

func main() {

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	if nil == ctx {
		fmt.Println("Get Background context fail, stop SMP!")
		return
	}
	defer cancel()

	var conf SysConfig
	initConfig(&conf)

	para_request := &db.DbConnPara{
		Host:     conf.cephManagerDbIp,
		Port:     conf.cephManagerDbPort,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "request"}

	cephManager := model.NewCephManager(conf.logDir, para_request)
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

	authManager := auth.NewAuthManager(conf.logDir, para_auth)
	if nil == authManager {
		fmt.Println("Start SMP AuthManager fail!")
		return
	}

	para_es := &db.DbConnPara{
		Host:     conf.elasticsearchManagerDbIp,
		Port:     conf.elasticsearchManagerDbPort,
		User:     "",
		Password: "",
		Dbname:   ""}

	monManager := model.NewMonitorManager(conf.logDir, para_es)
	if nil == monManager {
		fmt.Println("Start SMP monManager fail!")
		return
	}

	webService.SetupWebService(wg, ctx, cephManager, authManager, monManager)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	fmt.Println("Got signal:", s)
}
