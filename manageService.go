// manageService.go
package main

import (
	"net/http"
	"scheduler/common"
)

func addRequestType(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqType := req.FormValue("type")

	err := g_db.CreateNewRequestTable(reqType)
	var resp comm.CommonResponse

	if nil != err {
		resp.StateCode = comm.OP_ERROR
		resp.Msg = err.Error()
	} else {
		resp.StateCode = comm.OP_SUCCESS
		resp.Msg = "Create new request type success!"
	}

	resp.Send(rw)
}

func clean(rw http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	reqType := req.FormValue("type")

	err := g_db.RemoveRequestTable(reqType)
	var resp comm.CommonResponse

	if nil != err {
		resp.StateCode = comm.OP_ERROR
		resp.Msg = err.Error()
	} else {
		resp.StateCode = comm.OP_SUCCESS
		resp.Msg = "Remove request type" + reqType + " success"
	}

	resp.Send(rw)
}

func setupRequestService() {
	g_log.Info.Println("Listening Port 6667 for manage ...")
	http.HandleFunc("/addRequestType", addRequestType)
	http.HandleFunc("/clean", clean)
	http.ListenAndServe(":6667", nil)
}
