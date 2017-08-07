package webService

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/common"
	"strconv"
	"sync"
)

func fetchTask(rw http.ResponseWriter, req *http.Request) {
	reqType := req.FormValue("type")
	num, err := strconv.Atoi(req.FormValue("num"))
	if nil != err {
		g_reqHandler.InfoLog("Decode _num_ parameter fail,set to default num=1, ", err)
		num = 1
	}

	var resp comm.CommonResponse

	reqArr := g_reqHandler.GetUnprocessRequest(reqType, num)
	response := comm.RequestArray{Num: len(reqArr), RequestList: reqArr}

	b, err := json.Marshal(response)
	if nil == err {
		resp.StateCode = comm.OP_SUCCESS
		resp.Msg = string(b)
	} else {
		s := fmt.Sprintf("Encoding response fail, %s", err.Error())
		resp.StateCode = comm.OP_ERROR
		resp.Msg = s

		g_reqHandler.InfoLog(s)
	}
	fmt.Println("Get task!")
	resp.Send(rw)
}

func updateTask(rw http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	reqType := req.FormValue("type")

	decoder := json.NewDecoder(req.Body)
	var stateReport comm.ReqStateReport
	err := decoder.Decode(&stateReport)
	if nil != err {
		g_reqHandler.InfoLog("Decode request stateReport  fail, ", err)
		return
	}
	defer req.Body.Close()

	//print update task stateReport
	str, err := json.Marshal(&stateReport)
	if nil != err {
		g_reqHandler.InfoLog("Decode ReqStateReport fail", err)
		return
	}
	g_reqHandler.InfoLog("Update request:", reqType, string(str))

	if stateReport.State >= comm.REQUEST_STAT_TYPE_BUTT {
		g_reqHandler.InfoLog("Unknown State type", stateReport.State)
		return
	}

	err = g_reqHandler.UpdateRequestState(reqType, stateReport)
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
		fetchTask(rw, req)

	case "POST":
		fetchTask(rw, req)

	case "PUT":
		updateTask(rw, req)

	case "DELETE":
		s := "It's illegal to delete task!"
		g_reqHandler.InfoLog(s)
		rw.Write([]byte(s))

	default:
		g_reqHandler.InfoLog("Unknown request Method: ", req.Method)
	}

}

func setupWorkerService(wg *sync.WaitGroup, ctx context.Context) {

	server := &http.Server{Addr: _taskServAddr_, Handler: nil}

	http.HandleFunc("/task", taskOpDispatch)

	g_reqHandler.InfoLog("Listening Port 6668 for worker...")

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Stop Task Server.")
			server.Close()
			wg.Done()
		}
	}()

	server.ListenAndServe()
}
