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
var g_db db.DB

func addRequest(reqType string, request string) {

	id := fmt.Sprintf("%s", uuid.NewV4())
	t := time.Now()

	g_db.CreateNewRequest(id, t, reqType, request)
}

func getRequest(reqType string, num int) (reqNum int, reqs []comm.RequestWithUuid) {

	uuids := g_db.GetRequestInWaitingQueue(reqType, num)
	if len(uuids) == 0 {
		return
	}

	for _, id := range uuids {
		reqJson := g_db.GetSpecRequest(reqType, id)
		reqTemp := comm.RequestWithUuid{Id: id, Body: reqJson}
		reqs = append(reqs, reqTemp)
		reqNum += 1
	}
	return
}

func updateRequestState(reqType string, reqId string, workerId string, reqState comm.REQUEST_STATE_TYPE) error {

	currState, err := g_db.GetTaskState(reqType, reqId)
	if nil != err {
		g_log.Info.Println("Get task state fail, ", err)
		return err
	}

	if comm.REQUEST_IN_LINE == currState.State {
		g_db.RemoveRequestFromWaitingQueue(reqType, reqId)
	}

	currState.WorkerId = workerId
	currState.State = reqState
	currState.Timestamp[reqState] = time.Now()

	err = g_db.UpdateTaskState(nil, reqType, currState)
	return err
}

func main() {
	err := g_log.InitLogger("scheduler")
	if nil != err {
		log.Println(err.Error())
		return
	}
	g_db.InitDb("localhost", 6379, 10, "scheduler")
	defer func() {
		g_log.Info.Println("Scheduler stop!!")
	}()

	go setupRequestService()
	go setupWorkerService()

	for {
		time.Sleep(600 * time.Second)
	}
}
