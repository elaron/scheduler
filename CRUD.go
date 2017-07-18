package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
	"time"
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

func getSpecRequest(reqType string, reqId string) string {

	field := getReqTableName(reqType)
	resp := g_redisPool.Cmd("HGET", field, reqId)
	if nil != resp.Err {
		g_log.Info.Println("Get requst table fail, ", field, resp.Err)
		return "Get spec request fail"
	}

	body, err := resp.Str()
	if nil != err {
		s := fmt.Sprintf("Get request body fail, %s\n", err.Error())
		g_log.Info.Println(s)
		return s
	}
	return body
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
	updateTaskState(c, reqType, state)
}

func getTaskState(reqType string, reqId string) (RequestState, error) {

	field := getReqStateTableName(reqType)
	resp := g_redisPool.Cmd("HGET", field, reqId)
	if nil != resp.Err {
		return RequestState{}, resp.Err
	}

	str, err := resp.Str()
	if nil != err {
		return RequestState{}, err
	}

	var state RequestState
	err = json.Unmarshal([]byte(str), &state)
	if nil != err {
		return RequestState{}, err
	}

	return state, nil
}

func updateTaskState(c *redis.Client, reqType string, state RequestState) error {

	if nil == c {
		c, _ = g_redisPool.Get()
		defer g_redisPool.Put(c)
	}

	var msg string
	var err error
	var resp *redis.Resp
	var str []byte
	reqId := state.RequestId
	field := getReqStateTableName(reqType)

	str, err = json.Marshal(state)
	if nil != err {
		msg = fmt.Sprintf("Marsh request state fail, %s %s", reqId, err.Error())
		goto Fail
	}

	resp = c.Cmd("HSET", field, reqId, str)
	if nil != resp.Err {
		msg = fmt.Sprintf("Set request to state table fail, %s %s", reqId, resp.Err.Error())
		goto Fail
	}
	return nil

Fail:
	g_log.Info.Println(msg)
	return errors.New(msg)
}

func removeRequestFromWaitingQueue(reqType string, reqId string) {
	field := getReqWaitingQueueName(reqType)
	g_redisPool.Cmd("ZREM", field, reqId)
}
