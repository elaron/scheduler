package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

var g_log Log

func getReqTableName(reqType string) string {
	return fmt.Sprintf("request_%s", reqType)
}

func getReqWaitingQueueName(reqType string) string {
	return fmt.Sprintf("waiting_queue_%s", reqType)
}

func addRequestToTable(c redis.Conn, reqType string, reqId string, reqBody string) {

	field := getReqTableName(reqType)

	_, err := c.Do("HSET", field, reqId, reqBody)
	if nil != err {
		g_log.Info.Println("Add request to table fail, ", reqType, reqId, reqBody, err.Error())
		return
	}
}

func addReqToWatingQueue(c redis.Conn, reqType string, reqId string) {

	ts := time.Now().UnixNano()
	field := getReqWaitingQueueName(reqType)

	_, err := c.Do("ZADD", field, ts, reqId)
	if nil != err {
		g_log.Info.Println("Add request to waiting queue fail, ", reqType, reqId, err)
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

	addRequestToTable(c, reqType, id, request)
	addReqToWatingQueue(c, reqType, id)
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
		cleanRequestTable(c, reqType)
		c.Close()
	}()

	for i := 0; i < 5; i++ {
		addRequest(c, "100", "Elar")
	}
	printWaitingQueue(c, "100")
	printRequestTable(c, "100")
}
