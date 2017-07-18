package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"scheduler/common"
	"scheduler/log"
	"strconv"
	"time"
)

type DB struct {
	pool *pool.Pool
	log  logger.Log
}

func (db *DB) InitDb(address string, port int, poolSize int, prefix string) {
	addr := fmt.Sprintf("%s:%d", address, port)
	p, err := pool.New("tcp", addr, poolSize)
	if nil != err {
		db.log.Info.Println("Connect to redis fail, ", err)
		return
	}

	db.pool = p

	tag := fmt.Sprintf("%s-db-", prefix)
	db.log.InitLogger(tag)
}

func (db *DB) printResponse(r *redis.Resp) {
	db.log.Debug.Println(r.IsType(redis.Array))
	values, err := r.Map()
	if nil != err {
		db.log.Info.Println("Decode redis response fail ", err)
		return
	}

	for k, v := range values {
		db.log.Debug.Println(k, v)
	}
}

func (db *DB) PrintRequestStateTable(reqType string) {
	field := comm.GetReqStateTableName(reqType)
	resp := db.pool.Cmd("HGETALL", field)
	if nil != resp.Err {
		db.log.Info.Println("Get requst state table fail, ", field, resp.Err)
		return
	} else {
		db.log.Debug.Println("State table:")
		db.printResponse(resp)
	}
}

func (db *DB) PrintRequestTable(reqType string) {
	field := comm.GetReqTableName(reqType)
	resp := db.pool.Cmd("HGETALL", field)
	if nil != resp.Err {
		db.log.Info.Println("Get requst table fail, ", field, resp.Err)
		return
	} else {
		db.log.Debug.Println("Request table:")
		db.printResponse(resp)
	}
}

func (db *DB) GetSpecRequest(reqType string, reqId string) string {

	field := comm.GetReqTableName(reqType)
	resp := db.pool.Cmd("HGET", field, reqId)
	if nil != resp.Err {
		db.log.Info.Println("Get requst table fail, ", field, resp.Err)
		return "Get spec request fail"
	}

	body, err := resp.Str()
	if nil != err {
		s := fmt.Sprintf("Get request body fail, %s\n", err.Error())
		db.log.Info.Println(s)
		return s
	}
	return body
}

func (db *DB) PrintWaitingQueue(reqType string) {

	field := comm.GetReqWaitingQueueName(reqType)
	resp := db.pool.Cmd("ZRANGE", field, "0", "-1", "WITHSCORES")
	if nil != resp.Err {
		db.log.Info.Println("Get requst waiting queue fail, ", field, resp.Err)
		return
	} else {
		db.log.Debug.Println("Waiting queue:")
		db.printResponse(resp)
	}
}

func (db *DB) CleanRequestTable(reqType string) {
	field := comm.GetReqTableName(reqType)
	db.log.Debug.Println(">>>>>>>>>> Clean Request Table:")
	db.pool.Cmd("DEL", field)
}

func (db *DB) CleanRequestStateTable(reqType string) {
	field := comm.GetReqStateTableName(reqType)
	db.log.Debug.Println(">>>>>>>>>> Clean Request State Table:")
	db.pool.Cmd("DEL", field)
}

func (db *DB) CleanRequestWaitingQueue(reqType string) {
	wq := comm.GetReqWaitingQueueName(reqType)
	db.log.Debug.Println(">>>>>>>>>> Clean Request Waiting Queue:")
	resp := db.pool.Cmd("ZREMRANGEBYRANK", wq, "0", "-1")
	if nil != resp.Err {
		db.log.Info.Printf("Clean %s fail, %s\n", wq, resp.Err.Error())
		return
	}
}

func (db *DB) addRequestToTable(c *redis.Client, reqType string, reqId string, reqBody string) {

	field := comm.GetReqTableName(reqType)
	resp := c.Cmd("HSET", field, reqId, reqBody)
	if nil != resp.Err {
		db.log.Info.Println("Add request to table fail, ", reqType, reqId, reqBody, resp.Err.Error())
		return
	}
}

func (db *DB) addReqToWatingQueue(c *redis.Client, t time.Time, reqType string, reqId string) {

	ts := t.UnixNano()
	field := comm.GetReqWaitingQueueName(reqType)
	resp := c.Cmd("ZADD", field, ts, reqId)
	if nil != resp.Err {
		db.log.Info.Println("Add request to waiting queue fail, ", reqType, reqId, resp.Err)
		return
	}
}

func (db *DB) addReqToStateTable(c *redis.Client, t time.Time, reqType string, reqId string) {
	state := comm.RequestState{
		RequestId: reqId,
		State:     comm.REQUEST_IN_LINE,
	}
	state.Timestamp[comm.REQUEST_IN_LINE] = t
	db.UpdateTaskState(c, reqType, state)
}

func (db *DB) GetTaskState(reqType string, reqId string) (comm.RequestState, error) {

	field := comm.GetReqStateTableName(reqType)
	resp := db.pool.Cmd("HGET", field, reqId)
	if nil != resp.Err {
		return comm.RequestState{}, resp.Err
	}

	str, err := resp.Str()
	if nil != err {
		return comm.RequestState{}, err
	}

	var state comm.RequestState
	err = json.Unmarshal([]byte(str), &state)
	if nil != err {
		return comm.RequestState{}, err
	}

	return state, nil
}

func (db *DB) UpdateTaskState(c *redis.Client, reqType string, state comm.RequestState) error {

	if nil == c {
		c, _ = db.pool.Get()
		defer db.pool.Put(c)
	}

	var msg string
	var err error
	var resp *redis.Resp
	var str []byte
	reqId := state.RequestId
	field := comm.GetReqStateTableName(reqType)

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
	db.log.Info.Println(msg)
	return errors.New(msg)
}

func (db *DB) RemoveRequestFromWaitingQueue(reqType string, reqId string) {
	field := comm.GetReqWaitingQueueName(reqType)
	db.pool.Cmd("ZREM", field, reqId)
}

func (db *DB) CreateNewRequest(id string, t time.Time, reqType string, reqBody string) {
	c, err := db.pool.Get()
	if nil != err {
		db.log.Info.Println("Get connection from Redis Pool fail", err)
		return
	}
	defer db.pool.Put(c)

	resp := c.Cmd("MULTI")
	if nil != resp.Err {
		db.log.Info.Println(resp.Err)
		return
	}
	defer c.Cmd("EXEC")

	db.addRequestToTable(c, reqType, id, reqBody)
	db.addReqToWatingQueue(c, t, reqType, id)
	db.addReqToStateTable(c, t, reqType, id)
}

func (db *DB) GetRequestInWaitingQueue(reqType string, num int) []string {

	field := comm.GetReqWaitingQueueName(reqType)
	numStr := strconv.Itoa(num - 1)

	resp := db.pool.Cmd("ZRANGE", field, "0", numStr)
	if nil != resp.Err {
		db.log.Info.Println("Get requst waiting queue fail, ", field, resp.Err)
		return []string{}
	}

	uuids, err := resp.List()
	if nil != err {
		db.log.Info.Println("Decode waiting request queue fail, ", err)
		return []string{}
	}

	return uuids
}
