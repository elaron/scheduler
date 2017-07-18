package main

import (
	"fmt"
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

type ReqStateReport struct {
	RequestId string
	WorkerId  string
	State     REQUEST_STATE_TYPE
}

type RequestArray struct {
	Num         int
	RequestList []RequestWithUuid
}

type RequestWithUuid struct {
	Id   string
	Body string
}

func getReqTableName(reqType string) string {
	return fmt.Sprintf("request_%s", reqType)
}

func getReqWaitingQueueName(reqType string) string {
	return fmt.Sprintf("waiting_queue_%s", reqType)
}

func getReqStateTableName(reqType string) string {
	return fmt.Sprintf("req_state_%s", reqType)
}
