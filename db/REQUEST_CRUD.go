package db

import (
	"errors"
	"fmt"
	"scheduler/common"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//table req_state_reqType fields
const (
	RIT_FIELD_REQUEST_ID_        = "reqid"
	RIT_FIELD_SUBSCRIBE_         = "subscribe"
	RIT_FIELD_SUBSCRIBE_ADDRESS_ = "noticeaddress"
	RIT_FIELD_REQUEST_BODY_      = "reqbody"
)

//table req_state_reqType fields
const (
	RST_FIELD_REQUEST_ID_       = "reqid"
	RST_FIELD_WORKER_ID_        = "workerid"
	RST_FIELD_STATE_            = "state"
	RST_FIELD_CREATE_TIMESTAMP_ = "ts"
	RST_FIELD_UPDATE_TIMESTAMP_ = "updatets"
	RST_FIELD_RESPONSE_         = "resp"
)

func (p *Pgdb) CreateNewRequestTable(reqType string) error {

	var msg string
	db := p.db
	tx, err := db.Begin()
	if nil != err {
		msg = fmt.Sprintf("Begin transaction fail, %s", err.Error())
		return errors.New(msg)
	}
	defer tx.Commit()

	reqTableName := comm.GetReqTableName(reqType)
	reqStateTableName := comm.GetReqStateTableName(reqType)
	cmd := fmt.Sprintf("create table %s(%s varchar(64) PRIMARY KEY, %s boolean, %s varchar(1024), %s varchar(4096));create table %s(%s varchar(64) PRIMARY KEY, %s varchar(64), %s int, %s bigint, %s bigint, %s varchar(1024));",
		reqTableName,
		RIT_FIELD_REQUEST_ID_,
		RIT_FIELD_SUBSCRIBE_,
		RIT_FIELD_SUBSCRIBE_ADDRESS_,
		RIT_FIELD_REQUEST_BODY_,
		reqStateTableName,
		RST_FIELD_REQUEST_ID_,
		RST_FIELD_WORKER_ID_,
		RST_FIELD_STATE_,
		RST_FIELD_CREATE_TIMESTAMP_,
		RST_FIELD_UPDATE_TIMESTAMP_,
		RST_FIELD_RESPONSE_)

	_, err = tx.Exec(cmd)
	if nil != err {
		msg = fmt.Sprintf("Exec transaction fail, %s", err.Error())
		return errors.New(msg)
	}

	fmt.Println("Create new request table ", reqType)
	return nil
}

func (p *Pgdb) InsertNewRequest(reqType, reqId, reqBody, noticAddr string, subscribe bool) error {

	db := p.db
	reqestTable := comm.GetReqTableName(reqType)
	reqStateTable := comm.GetReqStateTableName(reqType)
	ts := time.Now().UnixNano()

	cmd := fmt.Sprintf("insert into %s(%s,%s,%s,%s) values('%s',%t,'%s','%s'); insert into %s(%s, %s, %s, %s, %s, %s) values('%s', '%s', %d, %d, %d, '%s');",
		reqestTable,
		RIT_FIELD_REQUEST_ID_, RIT_FIELD_SUBSCRIBE_, RIT_FIELD_SUBSCRIBE_ADDRESS_, RIT_FIELD_REQUEST_BODY_,
		reqId, subscribe, noticAddr, reqBody,
		reqStateTable,
		RST_FIELD_REQUEST_ID_, RST_FIELD_WORKER_ID_, RST_FIELD_STATE_, RST_FIELD_CREATE_TIMESTAMP_, RST_FIELD_UPDATE_TIMESTAMP_, RST_FIELD_RESPONSE_,
		reqId, "nobody", 0, ts, ts, "")

	_, err := db.Exec(cmd)
	if nil != err {
		return err
	}
	return nil
}

func (p *Pgdb) RemoveRequestTable(reqType string) error {
	db := p.db
	tx, err := db.Begin()
	if nil != err {
		return err
	}
	defer tx.Commit()

	reqTableName := comm.GetReqTableName(reqType)
	reqStateTable := comm.GetReqStateTableName(reqType)

	cmd := fmt.Sprintf("DROP TABLE IF EXISTS %s; DROP TABLE IF EXISTS %s;",
		reqTableName, reqStateTable)

	_, err = tx.Exec(cmd)
	if nil != err {
		return err
	}
	return nil
}

func (p *Pgdb) UpdateRequestState(reqType, reqId, workerid, resp string, reqState comm.REQUEST_STATE_TYPE) error {
	db := p.db

	reqStateTable := comm.GetReqStateTableName(reqType)
	cmd := fmt.Sprintf("update %s set %s = '%s', %s = %d, %s = '%s', %s = %d where %s = '%s';",
		reqStateTable,
		RST_FIELD_WORKER_ID_, workerid,
		RST_FIELD_STATE_, reqState,
		RST_FIELD_RESPONSE_, resp,
		RST_FIELD_UPDATE_TIMESTAMP_, time.Now().UnixNano(),
		RST_FIELD_REQUEST_ID_, reqId)

	_, err := db.Exec(cmd)
	if nil != err {
		return err
	}
	return nil
}

func (p *Pgdb) GetUnprocessRequest(reqType string, n int) (res []comm.RequestWithUuid, e error) {
	db := p.db
	reqTableName := comm.GetReqTableName(reqType)
	reqStateTable := comm.GetReqStateTableName(reqType)

	cmd := fmt.Sprintf("select * from %s where %s in (select %s from %s  where %s=0 order by ts limit %d);",
		reqTableName,
		RIT_FIELD_REQUEST_ID_,
		RST_FIELD_REQUEST_ID_,
		reqStateTable,
		RST_FIELD_STATE_,
		n)

	rows, err := db.Query(cmd)
	if nil != err {
		fmt.Println(err)
		e = err
		return
	}

	for rows.Next() {
		var reqId, reqBody, noticeAddr string
		var subscribe bool
		err = rows.Scan(&reqId, &subscribe, &noticeAddr, &reqBody)
		if nil != err {
			fmt.Println(err)
			e = err
			return
		}
		tmp := comm.RequestWithUuid{Id: reqId, Body: reqBody}
		res = append(res, tmp)
	}
	return
}

func (p *Pgdb) GetRequests(reqType string, reqIds []string) ([]comm.RequestInfo, error) {

	var riList []comm.RequestInfo

	cmd := fmt.Sprintf("select * from %s where %s in (%s);",
		comm.GetReqTableName(reqType),
		RIT_FIELD_REQUEST_ID_,
		strings.Join(reqIds, ","))

	rows, err := p.db.Query(cmd)
	if nil != err {
		return riList, err
	}

	for rows.Next() {
		var id, reqBody, noticeAddr string
		var subscribe bool
		err = rows.Scan(&id, &subscribe, &noticeAddr, &reqBody)
		if nil != err {
			return riList, err
		}

		ri := comm.RequestInfo{
			ReqType:   reqType,
			ReqId:     id,
			Subscribe: subscribe,
			SubAddr:   noticeAddr,
			ReqBody:   reqBody,
		}
		riList = append(riList, ri)
	}

	return riList, nil
}

func (p *Pgdb) GetRequestsState(reqType string, reqIds []string) ([]comm.RequestState, error) {

	var stateList []comm.RequestState
	reqIDsStr := fmt.Sprintf("'%s'", strings.Join(reqIds, "','"))

	cmd := fmt.Sprintf("select * from %s where %s in (%s);",
		comm.GetReqStateTableName(reqType),
		RST_FIELD_REQUEST_ID_,
		reqIDsStr)

	fmt.Print(cmd)

	rows, err := p.db.Query(cmd)
	if nil != err {
		return stateList, err
	}

	for rows.Next() {
		var rid, wid, resp string
		var cts, uts int64 //timestamp
		var state comm.REQUEST_STATE_TYPE

		err = rows.Scan(&rid, &wid, &state, &cts, &uts, &resp)
		if nil != err {
			return stateList, err
		}

		rs := comm.RequestState{
			RequestId:       rid,
			WorkerId:        wid,
			State:           state.String(),
			CreateTimestamp: time.Unix(0, cts),
			UpdateTimestamp: time.Unix(0, uts),
			Response:        resp,
		}
		stateList = append(stateList, rs)
	}

	return stateList, nil
}
