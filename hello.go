package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

var g_log Log

func addRequestToTable(c redis.Conn, reqType string, reqId string, reqBody string) {
	field := fmt.Sprintf("request_%s", reqType)
	_, err := c.Do("HSET", field, reqId, reqBody)
	if nil != err {

	}
}

func addRequest(c redis.Conn, request string) {

	ts := time.Now().UnixNano()
	id := uuid.NewV4()
	fmt.Println(ts, id)

	v, err := c.Do("ZADD", "request", ts, id)
	if nil != err {
		g_log.Info.Println("ZADD request fail", err)
		return
	}

	fmt.Println(v)
	values, err := redis.Values(c.Do("ZRANGE", "request", "0", "-1", "WITHSCORES"))
	if nil != err {
		g_log.Info.Println("Get requst list fail, ", err)
		return
	} else {
		fmt.Println("Request list:")
		for k, req := range values {
			if k%2 == 0 {
				fmt.Printf("id: %d %s ", k, req)
			} else {
				fmt.Printf("score: %s\n", req)
			}
		}
	}
}

func cleanRequestList(c redis.Conn, name string) {
	_, err := c.Do("ZREMRANGEBYRANK", name, "0", "-1")
	if nil != err {
		fmt.Printf("Clean %s fail, %s\n", name, err.Error())
		return
	}
}

func main() {

	err := InitLogger(&g_log, "scheduler")
	if nil != err {
		log.Println(err.Error())
		return
	}
	defer g_log.Info.Println("Scheduler stop!!")

	c, err := redis.Dial("tcp", "localhost:6379")
	if nil != err {
		g_log.Info.Println("Connect to redis fail, ", err)
		return
	}
	defer c.Close()
	cleanRequestList(c, "request")

	for i := 0; i < 5; i++ {
		addRequest(c, "Elar")
	}
}
