package model

import (
	"scheduler/db"
	"scheduler/log"
)

type CephManager struct {
	log *logger.Log
	db  *db.Pgdb

	requestHandler RequestHandler
}

func NewCephManager(dbParm *db.DbConnPara) *CephManager {
	manager := new(CephManager)

	manager.log = new(logger.Log)
	err := manager.log.InitLogger("CephManager")
	if nil != err {
		return nil
	}

	manager.db = db.NewPgDb(dbParm)
	if nil == manager.db {
		return nil
	}

	manager.requestHandler.Init(manager.log, manager.db)

	return manager
}

func (cm *CephManager) GetRequestHandler() *RequestHandler {
	return &cm.requestHandler
}
