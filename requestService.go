package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/common"
	"strconv"
	"strings"
)

func fetchRequest(rw http.ResponseWriter, req *http.Request) {
	reqType := req.FormValue("type")
	num, err := strconv.Atoi(req.FormValue("num"))
	if nil != err {
		g_log.Info.Println("Decode _num_ parameter fail, ", err)
		num = 1
	}
	g_log.Debug.Println("Get request:", reqType, num)

	reqNum, reqArr := getRequest(reqType, num)
	response := comm.RequestArray{Num: reqNum, RequestList: reqArr}
	b, err := json.Marshal(response)
	if nil != err {
		g_log.Info.Println("Encoding response fail", err)
		return
	}

	rw.Write(b)
}

func createRequest(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	reqType := req.FormValue("type")
	subscribe := req.FormValue("subscribe")
	noticeaddr := req.FormValue("noticeaddr")
	body := req.FormValue("body")

	sub := false
	if strings.ToUpper(subscribe) == "TRUE" {
		sub = true
	}

	fmt.Println(reqType, subscribe, noticeaddr, body)
	g_log.Debug.Println("Create request:", reqType, body)

	var resp comm.CommonResponse
	id, err := addRequest(reqType, noticeaddr, body, sub)

	if nil != err {
		resp.StateCode = comm.OP_ERROR
		resp.Msg = err.Error()
	} else {
		resp.StateCode = comm.OP_SUCCESS
		resp.Msg = id
	}

	resp.Send(rw)
}

func requestOpDispatch(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Action:", req.Method)
	switch req.Method {
	case "GET":
		fetchRequest(rw, req)

	case "POST":
		createRequest(rw, req)

	case "PUT":
		s := "It's illegal to update request!"
		g_log.Info.Println(s)
		rw.Write([]byte(s))

	case "DELETE":
		s := "It's illegal to update request!"
		g_log.Info.Println(s)
		rw.Write([]byte(s))

	default:
		g_log.Info.Println("Unknown request Method: ", req.Method)
	}

}

func setupManageService() {
	g_log.Info.Println("Listening Port 6666 for request...")
	http.HandleFunc("/request", requestOpDispatch)
	http.ListenAndServe(":6666", nil)
}
