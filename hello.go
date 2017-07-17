package main

import (
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
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
var g_redisPool *pool.Pool

func getReqTableName(reqType string) string {
	return fmt.Sprintf("request_%s", reqType)
}

func getReqWaitingQueueName(reqType string) string {
	return fmt.Sprintf("waiting_queue_%s", reqType)
}

func getReqStateTableName(reqType string) string {
	return fmt.Sprintf("req_state_%s", reqType)
}

func addRequestToTable(c *redis.Client, reqType string, reqId string, reqBody string) {

	field := getReqTableName(reqType)
	resp := c.Cmd("HSET", field, reqId, reqBody)
	if nil != resp.Err {
		g_log.Info.Println("Add request to table fail, ", reqType, reqId, reqBody, resp.Err.Error())
		return
	}
}

func addReqToWatingQueue(c *redis.Client, t time.Time, reqType string, reqId string) {

	ts := t.UnixNano()
	field := getReqWaitingQueueName(reqType)
	resp := c.Cmd("ZADD", field, ts, reqId)
	if nil != resp.Err {
		g_log.Info.Println("Add request to waiting queue fail, ", reqType, reqId, resp.Err)
		return
	}
}

func addReqToStateTable(c *redis.Client, t time.Time, reqType string, reqId string) {
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
	resp := c.Cmd("HSET", field, reqId, str)
	if nil != resp.Err {
		g_log.Info.Println("Add request to state table fail, ", reqId, resp.Err)
		return
	}
}

func addRequest(reqType string, request string) {

	id := fmt.Sprintf("%s", uuid.NewV4())
	c, err := g_redisPool.Get()
	if nil != err {
		g_log.Info.Println("Get connection from Redis Pool fail", err)
		return
	}

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

	setupService()
}
