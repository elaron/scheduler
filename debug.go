package main

import (
	"github.com/garyburd/redigo/redis"
)

func printValues(values []interface{}) {
	for _, v := range values {
		g_log.Debug.Printf("%s\n", v)
	}
}

func printRequestStateTable(c redis.Conn, reqType string) {
	field := getReqStateTableName(reqType)
	values, err := redis.Values(c.Do("HGETALL", field))
	if nil != err {
		g_log.Info.Println("Get requst state table fail, ", field, err)
		return
	} else {
		g_log.Debug.Println("State table:")
		printValues(values)
	}
}

func printRequestTable(c redis.Conn, reqType string) {
	field := getReqTableName(reqType)
	values, err := redis.Values(c.Do("HGETALL", field))
	if nil != err {
		g_log.Info.Println("Get requst table fail, ", field, err)
		return
	} else {
		g_log.Debug.Println("Request table:")
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
		g_log.Debug.Println("Waiting queue:")
		printValues(values)
	}
}

func cleanRequestTable(c redis.Conn, reqType string) {
	field := getReqTableName(reqType)
	cleanHTable(c, field)
}

func cleanRequestStateTable(c redis.Conn, reqType string) {
	field := getReqStateTableName(reqType)
	cleanHTable(c, field)
}

func cleanHTable(c redis.Conn, field string) {
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
