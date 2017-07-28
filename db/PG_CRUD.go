package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"scheduler/common"
	"time"
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
	cmd := fmt.Sprintf("create table %s(reqid varchar(64), reqbody varchar(4096));create table %s(reqid varchar(64), state int, ts bigint, resp varchar(1024));", reqTableName, reqStateTableName)

	_, err = tx.Exec(cmd)
	if nil != err {
		msg = fmt.Sprintf("Exec transaction fail, %s", err.Error())
		return errors.New(msg)
	}

	fmt.Println("Create new request table ", reqType)
	return nil
}

func (p *Pgdb) InsertNewRequest(reqType, reqId, reqBody string) error {
	db := p.db
	tx, err := db.Begin()
	if nil != err {
		return err
	}
	defer tx.Commit()

	reqestTable := comm.GetReqTableName(reqType)
	reqStateTable := comm.GetReqStateTableName(reqType)

	cmd := fmt.Sprintf("insert into %s(reqid,reqbody) values('%s','%s'); insert into %s(reqid, state, ts) values('%s', %d, %d);",
		reqestTable, reqId, reqBody,
		reqStateTable, reqId, 0, time.Now().UnixNano())

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

func (p *Pgdb) UpdateRequestState(reqType, reqId, resp string, reqState int) error {
	db := p.db

	reqStateTable := comm.GetReqStateTableName(reqType)
	cmd := fmt.Sprintf("update %s set state = %d, resp = '%s' where reqid = '%s';",
		reqStateTable, reqState, resp, reqId)
	fmt.Println("cmd>>", cmd)
	_, err := db.Exec(cmd)
	if nil != err {
		return err
	}
	return nil
}

func (p *Pgdb) GetUnprocessRequest(reqType string, n int) (res []string) {
	db := p.db
	reqTableName := comm.GetReqTableName(reqType)
	reqStateTable := comm.GetReqStateTableName(reqType)

	cmd := fmt.Sprintf("select * from %s where reqid in (select reqid from %s  where state=0 order by ts limit %d);",
		reqTableName, reqStateTable, n)

	rows, err := db.Query(cmd)
	if nil != err {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var reqId, reqBody string
		err = rows.Scan(&reqId, &reqBody)
		if nil != err {
			fmt.Println(err)
			return
		}
		res = append(res, reqId)
		res = append(res, reqBody)
		fmt.Println("reqid | reqbody ", reqId, reqBody)
	}
	return res
}
