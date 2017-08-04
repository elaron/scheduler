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
	log *logger.Log
	db  *db.Pgdb
}

func (rh *RequestHandler) Init(_log *logger.Log, _db *db.Pgdb) {
	rh.log = _log
	rh.db = _db
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

	stateList, err := rh.db.GetRequestsState(reqType, reqIDs)

	if nil != err {
		rh.log.Info.Println("Get request state fail, ", reqType, strings.Join(reqIDs, ","), err)
	}

	return stateList, err
}
