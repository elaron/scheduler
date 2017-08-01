package main

import (
	"fmt"
	"log"
	"scheduler/common"
	"scheduler/db"
	"scheduler/log"
	"time"

	"github.com/satori/go.uuid"
)

var g_log logger.Log
var g_db *db.Pgdb

func addRequest(reqType, noticeAddr, requestBody string, sub bool) (string, error) {

	id := fmt.Sprintf("%s", uuid.NewV4())
	err := g_db.InsertNewRequest(reqType, id, requestBody, noticeAddr, sub)
	if nil != err {
		g_log.Info.Println("Add request fail, ", err)
		return "", err
	}
	return id, nil
}

func getRequest(reqType string, num int) (reqNum int, reqs []comm.RequestWithUuid) {

	reqs, err := g_db.GetUnprocessRequest(reqType, num)
	if nil != err {
		g_log.Info.Println("Get unprocess request list fail, ", err)
		return
	}
	return
}

func updateRequestState(reqType string, sr comm.ReqStateReport) error {

	err := g_db.UpdateRequestState(reqType, sr.RequestId, sr.WorkerId, sr.Response, sr.State)
	return err
}

func main() {
	defer func() {
		g_log.Info.Println("Scheduler stop!!")
	}()

	err := g_log.InitLogger("scheduler")
	if nil != err {
		log.Println(err.Error())
		return
	}

	para := &db.DbConnPara{
		Host:     "192.168.56.132",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "request"}
	g_db = db.NewPgDb(para)
	if nil == g_db {
		g_log.Info.Println("Init postgresql fail!")
		return
	}

	go setupRequestService()
	go setupWorkerService()
	go setupManageService()

	for {
		time.Sleep(600 * time.Second)
	}
}
