package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"scheduler/common"
	"scheduler/db"
	"scheduler/log"
	"time"
)

var g_log logger.Log
var g_db *db.Pgdb

func addRequest(reqType string, request string) {

	id := fmt.Sprintf("%s", uuid.NewV4())
	err := g_db.InsertNewRequest(reqType, id, request, "", false)
	if nil != err {
		g_log.Info.Println("Add request fail, ", err)
		return
	}
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
		Host:     "127.0.0.1",
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

	for {
		time.Sleep(600 * time.Second)
	}
}
