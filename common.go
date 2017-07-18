package main

import (
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
