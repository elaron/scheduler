package comm

import (
	"fmt"
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
	RequestId       string
	WorkerId        string
	State           REQUEST_STATE_TYPE
	CreateTimestamp int64
	UpdateTimestamp int64
	Response        string
}

type ReqStateReport struct {
	RequestId string
	WorkerId  string
	State     REQUEST_STATE_TYPE
	Response  string
}

type RequestArray struct {
	Num         int
	RequestList []RequestWithUuid
}

type RequestWithUuid struct {
	Id   string
	Body string
}

func GetReqTableName(reqType string) string {
	return fmt.Sprintf("request_%s", reqType)
}

func GetReqWaitingQueueName(reqType string) string {
	return fmt.Sprintf("waiting_queue_%s", reqType)
}

func GetReqStateTableName(reqType string) string {
	return fmt.Sprintf("req_state_%s", reqType)
}

func GetUserTable(username string) string {
	return fmt.Sprintf("user_%s", username)
}
