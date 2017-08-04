package model

import (
	"fmt"
	"scheduler/common"
	"scheduler/db"
	"scheduler/log"
	"strings"

	"github.com/satori/go.uuid"
)

type RequestHandler struct {
	log      *logger.Log
	db       *db.Pgdb
	statStr  map[comm.REQUEST_STATE_TYPE]string
	statEnum map[string]comm.REQUEST_STATE_TYPE
}

func (rh *RequestHandler) Init(_log *logger.Log, _db *db.Pgdb) {
	rh.log = _log
	rh.db = _db

	rh.statStr = make(map[comm.REQUEST_STATE_TYPE]string)
	rh.statEnum = make(map[string]comm.REQUEST_STATE_TYPE)

	rh.statStr[comm.REQUEST_IN_LINE] = "IN_LINE"
	rh.statStr[comm.REQUEST_BEGIN_PRO] = "BEGIN_PRO"
	rh.statStr[comm.REQUEST_PAUSE] = "PAUSE"
	rh.statStr[comm.REQUEST_FINISH] = "FINISH"
	rh.statStr[comm.REQUEST_CANCEL] = "CANCEL"

	rh.statEnum["IN_LINE"] = comm.REQUEST_IN_LINE
	rh.statEnum["BEGIN_PRO"] = comm.REQUEST_BEGIN_PRO
	rh.statEnum["PAUSE"] = comm.REQUEST_PAUSE
	rh.statEnum["FINISH"] = comm.REQUEST_FINISH
	rh.statEnum["CANCEL"] = comm.REQUEST_CANCEL
}

func (rh *RequestHandler) StateStr(state comm.REQUEST_STATE_TYPE) string {
	if state >= comm.REQUEST_STAT_TYPE_BUTT {
		return "illegalType"
	}

	return rh.statStr[state]
}

func (rh *RequestHandler) InfoLog(v ...interface{}) {
	rh.log.Info.Println(v)
}

func (rh *RequestHandler) DebugLog(v ...interface{}) {
	rh.log.Debug.Println(v)
}

func (rh *RequestHandler) CreateNewRequestType(reqType string) error {
	rh.log.Info.Println("Create New Request Type ", reqType)
	return rh.db.CreateNewRequestTable(reqType)
}

func (rh *RequestHandler) DeleteRequestType(reqType string) error {
	rh.log.Info.Println("Delete Request Type ", reqType)
	return rh.db.RemoveRequestTable(reqType)
}

func (rh *RequestHandler) AddNewRequest(reqType, noticeAddr, requestBody string, sub bool) (string, error) {
	id := fmt.Sprintf("%s", uuid.NewV4())
	err := rh.db.InsertNewRequest(reqType, id, requestBody, noticeAddr, sub)
	if nil != err {
		rh.log.Info.Println("Add request fail, ", err, reqType, requestBody)
		return "", err
	}
	return id, nil
}

func (rh *RequestHandler) GetUnprocessRequest(reqType string, num int) (reqs []comm.RequestWithUuid) {

	reqs, err := rh.db.GetUnprocessRequest(reqType, num)
	if nil != err {
		rh.log.Info.Println("Get unprocess request list fail, ", err)
		return
	}
	return
}

func (rh *RequestHandler) UpdateRequestState(reqType string, sr comm.ReqStateReport) error {

	err := rh.db.UpdateRequestState(reqType, sr.RequestId, sr.WorkerId, sr.Response, sr.State)

	if nil == err {
		rh.log.Info.Println("Update Request State success, ", reqType, sr.RequestId, sr.WorkerId, sr.State, sr.Response)
	} else {
		rh.log.Info.Println("Update Request State fail, ", reqType, sr.RequestId, sr.WorkerId, sr.State, sr.Response)
	}

	return err
}

func (rh *RequestHandler) GetRequestsState(reqType string, reqIDs []string) ([]comm.RequestState, error) {

	stateList, err := rh.db.GetRequestsState(reqType, reqIDs, rh.StateStr)

	if nil != err {
		rh.log.Info.Println("Get request state fail, ", reqType, strings.Join(reqIDs, ","), err)
	}

	return stateList, err
}
