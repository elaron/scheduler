package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
)

var g_log Log
var g_redisPool *pool.Pool

func addRequest(reqType string, request string) {

	id := fmt.Sprintf("%s", uuid.NewV4())
	c, err := g_redisPool.Get()
	if nil != err {
		g_log.Info.Println("Get connection from Redis Pool fail", err)
		return
	}
	defer g_redisPool.Put(c)

	resp := c.Cmd("MULTI")
	if nil != resp.Err {
		g_log.Info.Println(resp.Err)
		return
	}
	defer c.Cmd("EXEC")

	t := time.Now()
	addRequestToTable(c, reqType, id, request)
	addReqToWatingQueue(c, t, reqType, id)
	addReqToStateTable(c, t, reqType, id)
}

func getRequest(reqType string, num int) (reqNum int, reqs []RequestWithUuid) {

	field := getReqWaitingQueueName(reqType)
	numStr := strconv.Itoa(num - 1)

	resp := g_redisPool.Cmd("ZRANGE", field, "0", numStr)
	if nil != resp.Err {
		g_log.Info.Println("Get requst waiting queue fail, ", field, resp.Err)
		return

	} else {
		uuids, err := resp.List()
		if nil != err {
			g_log.Info.Println("Decode waiting request queue fail, ", err)
			return
		}

		for _, id := range uuids {
			reqJson := getSpecRequest(reqType, id)
			reqTemp := RequestWithUuid{Id: id, Body: reqJson}
			reqs = append(reqs, reqTemp)
			reqNum += 1
		}
	}
	return
}

func updateRequestState(reqType string, reqId string, workerId string, reqState REQUEST_STATE_TYPE) error {

	currState, err := getTaskState(reqType, reqId)
	if nil != err {
		g_log.Info.Println("Get task state fail, ", err)
		return err
	}

	currState.WorkerId = workerId
	currState.State = reqState
	currState.Timestamp[reqState] = time.Now()

	err = updateTaskState(nil, reqType, currState)
	return err
}

func main() {
	err := InitLogger(&g_log, "scheduler")
	if nil != err {
		log.Println(err.Error())
		return
	}

	p, err := pool.New("tcp", "localhost:6379", 10)
	if nil != err {
		g_log.Info.Println("Connect to redis fail, ", err)
		return
	}
	g_redisPool = p

	defer func() {
		g_log.Info.Println("Scheduler stop!!")
	}()

	go setupRequestService()
	go setupWorkerService()

	for {
		time.Sleep(600 * time.Second)
	}
}
