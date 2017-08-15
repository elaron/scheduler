package webService

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/common"
	"strings"
	"sync"
)

func setupRequestService(wg *sync.WaitGroup, ctx context.Context) {

	server := &http.Server{Addr: _requestServAddr_, Handler: nil}

	http.HandleFunc("/request", requestOpDispatch)
	http.HandleFunc("/requestState", requestStateOpDispatch)

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Stop Request Server.")
			server.Close()
			wg.Done()
		}
	}()
	g_reqHandler.InfoLog("Listening for request...on port ", _requestServAddr_)
	server.ListenAndServe()
}

func requestOpDispatch(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Action:", req.Method)
	switch req.Method {

	case "POST":
		createRequest(rw, req)

	default:
		s := fmt.Sprintf("It's illegal to %s request!", req.Method)
		g_reqHandler.InfoLog(s)
		rw.Write([]byte(s))
	}
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

func getRequestState(rw http.ResponseWriter, req *http.Request) {

	reqType := req.FormValue("type")
	reqIds := req.FormValue("reqIDs")

	g_reqHandler.InfoLog("Get request:", reqType, reqIds)

	if len(reqType) == 0 || len(reqIds) == 0 {
		rw.Write([]byte("_type_ or _reqIDs_ field illegal"))
		return
	}

	var resp comm.CommonResponse
	reqIdList := strings.Split(reqIds, ",")

	stateList, err := g_reqHandler.GetRequestsState(reqType, reqIdList)

	if nil != err {
		resp.StateCode = comm.OP_ERROR
		resp.Msg = err.Error()
	} else {
		response := comm.RequestStateArray{Num: len(stateList), StateList: stateList}
		b, _ := json.Marshal(response)
		resp.StateCode = comm.OP_SUCCESS
		resp.Msg = string(b)
	}

	resp.Send(rw)
}

func requestStateOpDispatch(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Action:", req.Method)
	switch req.Method {

	case "GET":
		getRequestState(rw, req)

	case "POST":
		getRequestState(rw, req)

	default:
		s := fmt.Sprintf("It's illegal to %s request state!", req.Method)
		g_reqHandler.InfoLog(s)
		rw.Write([]byte(s))
	}
}
