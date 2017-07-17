package main

import (
	"github.com/mediocregopher/radix.v2/redis"
)

func printResponse(r *redis.Resp) {
	g_log.Debug.Println(r.IsType(redis.Array))
	values, err := r.Map()
	if nil != err {
		g_log.Info.Println("Decode redis response fail ", err)
		return
	}

	for k, v := range values {
		g_log.Debug.Println(k, v)
	}
}

func printRequestStateTable(reqType string) {
	field := getReqStateTableName(reqType)
	resp := g_redisPool.Cmd("HGETALL", field)
	if nil != resp.Err {
		g_log.Info.Println("Get requst state table fail, ", field, resp.Err)
		return
	} else {
		g_log.Debug.Println("State table:")
		printResponse(resp)
	}
}

func printRequestTable(reqType string) {
	field := getReqTableName(reqType)
	resp := g_redisPool.Cmd("HGETALL", field)
	if nil != resp.Err {
		g_log.Info.Println("Get requst table fail, ", field, resp.Err)
		return
	} else {
		g_log.Debug.Println("Request table:")
		printResponse(resp)
	}
}

func printWaitingQueue(reqType string) {

	field := getReqWaitingQueueName(reqType)
	resp := g_redisPool.Cmd("ZRANGE", field, "0", "-1", "WITHSCORES")
	if nil != resp.Err {
		g_log.Info.Println("Get requst waiting queue fail, ", field, resp.Err)
		return
	} else {
		g_log.Debug.Println("Waiting queue:")
		printResponse(resp)
	}
}

func cleanRequestTable(reqType string) {
	field := getReqTableName(reqType)
	g_log.Debug.Println(">>>>>>>>>> Clean Request Table:")
	cleanHTable(field)
}

func cleanRequestStateTable(reqType string) {
	field := getReqStateTableName(reqType)
	g_log.Debug.Println(">>>>>>>>>> Clean Request State Table:")
	cleanHTable(field)
}

func cleanHTable(field string) {
	resp := g_redisPool.Cmd("HGETALL", field)
	if resp.Err != nil {
		g_log.Info.Printf("Clean %s fail, %s\n", field, resp.Err.Error())
		return
	}

	resp = g_redisPool.Cmd("HGETALL", field)
	if nil != resp.Err {
		g_log.Info.Println("Get requst table fail, ", field, resp.Err)
		return
	} else {
		if false == resp.IsType(redis.Array) {
			g_log.Info.Println("Decode fail", field)
			return
		}

		values, err := resp.Map()
		if nil != err {
			g_log.Info.Println("Decode Htable response fail", field)
			return
		}
		for k, v := range values {
			g_redisPool.Cmd("HDEL", field, k)
			g_log.Debug.Println("Removing ", k, v)
		}
	}
}

func cleanRequestWaitingQueue(reqType string) {
	wq := getReqWaitingQueueName(reqType)
	g_log.Debug.Println(">>>>>>>>>> Clean Request Waiting Queue:")
	resp := g_redisPool.Cmd("ZREMRANGEBYRANK", wq, "0", "-1")
	if nil != resp.Err {
		g_log.Info.Printf("Clean %s fail, %s\n", wq, resp.Err.Error())
		return
	}
}
