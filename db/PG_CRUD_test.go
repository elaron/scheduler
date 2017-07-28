package db

import (
	"fmt"
	"testing"
)

func Test_CreateRequestTable_1(t *testing.T) {
	para := &DbConnPara{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "request"}

	p := NewPgDb(para)
	reqType := "100"
	err := p.RemoveRequestTable(reqType)
	if nil != err {
		t.Error("CASE: Test_CreateRequestTable_1 Remove table FAIL!", err)
		return
	}

	err = p.CreateNewRequestTable(reqType)
	if nil != err {
		t.Error("[CASE: Test_CreateRequestTable_1 create table FAIL!] ", err)
		return
	}
	for i := 0; i < 5; i++ {
		reqId := fmt.Sprintf("abcd%d", i)
		reqBody := fmt.Sprintf("abcd%d_body", i)
		err = p.InsertNewRequest(reqType, reqId, reqBody)
		if nil != err {
			t.Error("CASE: Test_CreateRequestTable_1 insert new requestFAIL!", err)
			return
		}
	}

	t.Log("[CASE Test_CreateRequestTable_1 success]")
}
