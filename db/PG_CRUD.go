package db

import (
	"database/sql"
	"errors"
	"fmt"
	"scheduler/common"
	"time"

	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "127.0.0.1"
	DB_PORT     = 5432
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "request"
)

type Pgdb struct {
	db *sql.DB
}

type DbConnPara struct {
	Host     string
	Port     int32
	User     string
	Password string
	Dbname   string
}

func NewPgDb(para *DbConnPara) *Pgdb {

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		para.Host, para.Port, para.User, para.Password, para.Dbname)

	db, err := sql.Open("postgres", dbInfo)
	if nil != err {
		fmt.Println("Open fail, ", err)
		return nil
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return &Pgdb{db: db}
}

func (p *Pgdb) Close() {
	p.db.Close()
}

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
	cmd := fmt.Sprintf("create table %s(reqid varchar(64), subscribe boolean, noticeaddress varchar(1024), reqbody varchar(4096));create table %s(reqid varchar(64), workerid varchar(64), state int, ts bigint, updatets bigint, resp varchar(1024));", reqTableName, reqStateTableName)

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
	tx, err := db.Begin()
	if nil != err {
		return err
	}
	defer tx.Commit()

	reqestTable := comm.GetReqTableName(reqType)
	reqStateTable := comm.GetReqStateTableName(reqType)
	ts := time.Now().UnixNano()

	cmd := fmt.Sprintf("insert into %s(reqid,subscribe,noticeaddress,reqbody) values('%s',%t,'%s','%s'); insert into %s(reqid, state, ts, updatets) values('%s', %d, %d, %d);",
		reqestTable, reqId, subscribe, noticAddr, reqBody,
		reqStateTable, reqId, 0, ts, ts)

	_, err = tx.Exec(cmd)
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
	cmd := fmt.Sprintf("update %s set workerid = '%s', state = %d, resp = '%s', updatets = %d where reqid = '%s';",
		reqStateTable, workerid, reqState, resp, time.Now().UnixNano(), reqId)
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

	cmd := fmt.Sprintf("select * from %s where reqid in (select reqid from %s  where state=0 order by ts limit %d);",
		reqTableName, reqStateTable, n)

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
		fmt.Println("reqid | reqbody ", reqId, reqBody)
	}
	return
}
