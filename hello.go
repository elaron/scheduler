package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

const (
	REQUEST_IN_LINE = iota
	REQUEST_BEGIN_PRO
	REQUEST_PAUSE
	REQUEST_FINISH
	REQUEST_CANCEL
	REQUEST_STAT_TYPE_BUTT
)

type REQUEST_STATE_TYPE int32

type RequestState struct {
	RequestId string
	WorkerId  string
	State     REQUEST_STATE_TYPE
	Timestamp [REQUEST_STAT_TYPE_BUTT]time.Time
}

var g_log Log

func getReqTableName(reqType string) string {
	return fmt.Sprintf("request_%s", reqType)
}

func getReqWaitingQueueName(reqType string) string {
	return fmt.Sprintf("waiting_queue_%s", reqType)
}

func getReqStateTableName(reqType string) string {
	return fmt.Sprintf("req_state_%s", reqType)
}

func addRequestToTable(c redis.Conn, reqType string, reqId string, reqBody string) {

	field := getReqTableName(reqType)
	_, err := c.Do("HSET", field, reqId, reqBody)
	if nil != err {
		g_log.Info.Println("Add request to table fail, ", reqType, reqId, reqBody, err.Error())
		return
	}
}

func addReqToWatingQueue(c redis.Conn, t time.Time, reqType string, reqId string) {

	ts := t.UnixNano()
	field := getReqWaitingQueueName(reqType)
	_, err := c.Do("ZADD", field, ts, reqId)
	if nil != err {
		g_log.Info.Println("Add request to waiting queue fail, ", reqType, reqId, err)
		return
	}
}

func addReqToStateTable(c redis.Conn, t time.Time, reqType string, reqId string) {
	state := RequestState{
		RequestId: reqId,
		State:     REQUEST_IN_LINE,
	}
	state.Timestamp[REQUEST_IN_LINE] = t

	str, err := json.Marshal(state)
	if nil != err {
		g_log.Info.Println("Marsh request state fail, ", reqId, err)
		return
	}

	field := getReqStateTableName(reqType)
	_, err = c.Do("HSET", field, reqId, str)
	if nil != err {
		g_log.Info.Println("Add request to state table fail, ", reqId, err)
		return
	}
}

func addRequest(c redis.Conn, reqType string, request string) {

	id := fmt.Sprintf("%s", uuid.NewV4())
	//start a transaction
	_, err := c.Do("MULTI")
	if nil != err {
		g_log.Info.Println(err)
	}
	defer c.Do("EXEC")

	t := time.Now()
	addRequestToTable(c, reqType, id, request)
	addReqToWatingQueue(c, t, reqType, id)
	addReqToStateTable(c, t, reqType, id)
}

func main() {
	reqType := "100"
	err := InitLogger(&g_log, "scheduler")
	if nil != err {
		log.Println(err.Error())
		return
	}

	c, err := redis.Dial("tcp", "localhost:6379")
	if nil != err {
		g_log.Info.Println("Connect to redis fail, ", err)
		return
	}

	defer func() {
		g_log.Info.Println("Scheduler stop!!")
		cleanRequestWaitingQueue(c, reqType)
		cleanRequestStateTable(c, reqType)
		cleanRequestTable(c, reqType)
		c.Close()
	}()

	for i := 0; i < 5; i++ {
		addRequest(c, reqType, "Elar")
	}
	printWaitingQueue(c, reqType)
	printRequestStateTable(c, reqType)
	printRequestTable(c, reqType)
	setupService()
}
