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
	for i := 0; i < 10; i++ {
		reqId := fmt.Sprintf("abcd%d", i)
		reqBody := fmt.Sprintf("abcd%d_body", i)
		err = p.InsertNewRequest(reqType, reqId, reqBody)
		if nil != err {
			t.Error("CASE: Test_CreateRequestTable_1 insert new requestFAIL!", err)
			return
		}

		if i%2 == 0 {
			err = p.UpdateRequestState(reqType, reqId, "some response", 2)
			if nil != err {
				t.Error("CASE: Test_CreateRequestTable_1 update  requestFAIL!", err)
				return
			}

		}
	}
	arr := p.GetUnprocessRequest(reqType, 3)
	if 3 != len(arr) {
		t.Error("CASE: Test_CreateRequestTable_1 get  requestFAIL!")
		return
	}
	t.Log(arr)

	t.Log("[CASE Test_CreateRequestTable_1 success]")
}
