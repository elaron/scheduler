package webService

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
		g_reqHandler.InfoLog("Decode _num_ parameter fail, ", err)
		num = 1
	}
	g_reqHandler.InfoLog("Get request:", reqType, num)

	reqArr := g_reqHandler.GetUnprocessRequest(reqType, num)
	response := comm.RequestArray{Num: len(reqArr), RequestList: reqArr}
	b, err := json.Marshal(response)
	if nil != err {
		g_reqHandler.InfoLog("Encoding response fail", err)
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
	g_reqHandler.InfoLog("Create request:", reqType, body)

	var resp comm.CommonResponse
	id, err := g_reqHandler.AddNewRequest(reqType, noticeaddr, body, sub)

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
		g_reqHandler.InfoLog(s)
		rw.Write([]byte(s))

	case "DELETE":
		s := "It's illegal to update request!"
		g_reqHandler.InfoLog(s)
		rw.Write([]byte(s))

	default:
		g_reqHandler.InfoLog("Unknown request Method: ", req.Method)
	}

}

func setupManageService() {
	g_reqHandler.InfoLog("Listening Port 6666 for request...")
	http.HandleFunc("/request", requestOpDispatch)
	http.ListenAndServe(":6666", nil)
}
