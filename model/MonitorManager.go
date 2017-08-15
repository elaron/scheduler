package model

import (
	"fmt"
	"scheduler/db"
	"scheduler/log"
)

type MonitorManager struct {
	log *logger.Log
	db  *db.EsDb
}

func NewMonitorManager(logDir string, dbParm *db.DbConnPara) *MonitorManager {
	auth := new(MonitorManager)

	auth.log = new(logger.Log)
	err := auth.log.InitLogger(logDir, "Auth")
	if nil != err {
		fmt.Println(err)
		return nil
	}

	auth.db = db.NewEsDb(dbParm)
	if nil == auth.db {
		return nil
	}

	return auth
}
