package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func fetchTask(rw http.ResponseWriter, req *http.Request) {
	reqType := req.FormValue("type")
	num, err := strconv.Atoi(req.FormValue("num"))
	if nil != err {
		g_log.Info.Println("Decode _num_ parameter fail, ", err)
		return
	}

	reqNum, reqArr := getRequest(reqType, num)
	response := RequestArray{Num: reqNum, RequestList: reqArr}
	b, err := json.Marshal(response)
	if nil != err {
		g_log.Info.Println("Encoding response fail", err)
		return
	}

	rw.Write(b)
}

func updateTask(rw http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	reqType := req.FormValue("type")

	decoder := json.NewDecoder(req.Body)
	var stateReport ReqStateReport
	err := decoder.Decode(&stateReport)
	if nil != err {
		g_log.Info.Println("Decode request stateReport  fail, ", err)
		return
	}
	defer req.Body.Close()

	//print update task stateReport
	str, err := json.Marshal(&stateReport)
	if nil != err {
		g_log.Info.Println("Decode ReqStateReport fail", err)
		return
	}
	g_log.Debug.Println("Update request:", reqType, string(str))

	if stateReport.State >= REQUEST_STAT_TYPE_BUTT {
		g_log.Info.Println("Unknown State type", stateReport.State)
		return
	}

	currState, err := getTaskState(reqType, stateReport.RequestId)
	if nil != err {
		g_log.Info.Println("Get task state fail, ", err)
		return
	}

	currState.WorkerId = stateReport.WorkerId
	currState.State = stateReport.State
	currState.Timestamp[stateReport.State] = time.Now()

	err = updateTaskState(nil, reqType, currState)
	if nil != err {
		rw.Write([]byte(err.Error()))
	}
	s := fmt.Sprintf("Update task %s %d success!", stateReport.RequestId, stateReport.State)
	rw.Write([]byte(s))
}

func taskOpDispatch(rw http.ResponseWriter, req *http.Request) {

	fmt.Println("Action:", req.Method)

	switch req.Method {
	case "GET":
		fetchRequest(rw, req)

	case "POST":
		s := "It's illegal to create task!"
		g_log.Info.Println(s)
		rw.Write([]byte(s))

	case "PUT":
		updateTask(rw, req)

	case "DELETE":
		s := "It's illegal to delete task!"
		g_log.Info.Println(s)
		rw.Write([]byte(s))

	default:
		g_log.Info.Println("Unknown request Method: ", req.Method)
	}

}

func setupWorkerService() {
	g_log.Info.Println("Listening Port 2345 for worker...")
	http.HandleFunc("/task", taskOpDispatch)
	http.ListenAndServe(":2345", nil)
}
