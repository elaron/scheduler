package db

import (
	"fmt"
	"testing"
)

func Test_CreateRequestTable_1(t *testing.T) {
	para := &DbConnPara{
		Host:     "192.168.65.132",
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
	for i := 0; i < 10; i++ {
		reqId := fmt.Sprintf("abcd%d", i)
		reqBody := fmt.Sprintf("abcd%d_body", i)

		if i%3 == 0 {
			err = p.InsertNewRequest(reqType, reqId, reqBody, "notice me here", true)
		} else {
			err = p.InsertNewRequest(reqType, reqId, reqBody, "", false)
		}
		if nil != err {
			t.Error("CASE: Test_CreateRequestTable_1 insert new requestFAIL!", err)
			return
		}

		workderid := "wokerNo1"
		if i%2 == 0 {
			err = p.UpdateRequestState(reqType, reqId, workderid, "some response", 2)
			if nil != err {
				t.Error("CASE: Test_CreateRequestTable_1 update  requestFAIL!", err)
				return
			}

		}
	}
	arr, err := p.GetUnprocessRequest(reqType, 3)
	if 3 != len(arr) || nil != err {
		t.Error("CASE: Test_CreateRequestTable_1 get  requestFAIL!")
		return
	}
	t.Log(arr)

	t.Log("[CASE Test_CreateRequestTable_1 success]")
}
