package comm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RequestInfo struct {
	ReqType   string
	ReqId     string
	Subscribe bool
	SubAddr   string
	ReqBody   string
}

const (
	REQUEST_IN_LINE = iota
	REQUEST_BEGIN_PRO
	REQUEST_PAUSE
	REQUEST_FINISH
	REQUEST_CANCEL
	REQUEST_STAT_TYPE_BUTT
)

type StateEnum2Str func(REQUEST_STATE_TYPE) string
type Str2StateEnum func(string) REQUEST_STATE_TYPE

type REQUEST_STATE_TYPE int32

type RequestState struct {
	RequestId       string
	WorkerId        string
	State           string
	CreateTimestamp time.Time
	UpdateTimestamp time.Time
	Response        string
}

type ReqStateReport struct {
	RequestId string
	WorkerId  string
	State     REQUEST_STATE_TYPE
	Response  string
}

type RequestStateArray struct {
	Num       int
	StateList []RequestState
}

type RequestArray struct {
	Num         int
	RequestList []RequestWithUuid
}

type RequestWithUuid struct {
	Id   string
	Body string
}

const (
	OP_SUCCESS = 100
	OP_ERROR   = 900
)

type CommonResponse struct {
	StateCode int
	Msg       string
}

func (cr *CommonResponse) Send(rw http.ResponseWriter) error {
	b, err := json.Marshal(cr)
	rw.Write(b)
	return err
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
