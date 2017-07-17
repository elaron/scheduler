package main

import (
	"github.com/garyburd/redigo/redis"
)

func printValues(values []interface{}) {
	g_log.Debug.Println("list:")
	for _, v := range values {
		g_log.Debug.Printf("%s\n", v)
	}
}

func printRequestTable(c redis.Conn, reqType string) {
	field := getReqTableName(reqType)
	values, err := redis.Values(c.Do("HGETALL", field))
	if nil != err {
		g_log.Info.Println("Get requst table fail, ", field, err)
		return
	} else {
		printValues(values)
	}
}

func printWaitingQueue(c redis.Conn, reqType string) {

	field := getReqWaitingQueueName(reqType)
	values, err := redis.Values(c.Do("ZRANGE", field, "0", "-1", "WITHSCORES"))
	if nil != err {
		g_log.Info.Println("Get requst waiting queue fail, ", field, err)
		return
	} else {
		printValues(values)
	}
}

func cleanRequestTable(c redis.Conn, reqType string) {
	field := getReqTableName(reqType)
	_, err := c.Do("HGETALL", field)
	if nil != err {
		g_log.Info.Printf("Clean %s fail, %s\n", field, err.Error())
		return
	}

	values, err := redis.Values(c.Do("HGETALL", field))
	if nil != err {
		g_log.Info.Println("Get requst table fail, ", field, err)
		return
	} else {
		for k, v := range values {
			if k%2 == 0 {
				c.Do("HDEL", field, v)
			}
		}
	}
}

func cleanRequestWaitingQueue(c redis.Conn, reqType string) {
	wq := getReqWaitingQueueName(reqType)
	_, err := c.Do("ZREMRANGEBYRANK", wq, "0", "-1")
	if nil != err {
		g_log.Info.Printf("Clean %s fail, %s\n", wq, err.Error())
		return
	}
}
