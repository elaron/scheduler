package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ReqStateReport struct {
	RequestId string
	WorkerId  string
	State     REQUEST_STATE_TYPE
}

type RequestArray struct {
	Num         int
	RequestList []RequestWithUuid
}

func fetchRequest(rw http.ResponseWriter, req *http.Request) {
	reqType := req.FormValue("type")
	num, err := strconv.Atoi(req.FormValue("num"))
	if nil != err {
		g_log.Info.Println("Decode _num_ parameter fail, ", err)
		return
	}
	g_log.Debug.Println("Get request:", reqType, num)

	printWaitingQueue(reqType)

	reqNum, reqArr := getRequest(reqType, num)
	response := RequestArray{Num: reqNum, RequestList: reqArr}
	b, err := json.Marshal(response)
	if nil != err {
		g_log.Info.Println("Encoding response fail", err)
		return
	}

	rw.Write(b)
}

func updateRequest(rw http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	reqType := req.FormValue("type")

	decoder := json.NewDecoder(req.Body)
	var state ReqStateReport
	err := decoder.Decode(&state)
	if nil != err {
		g_log.Info.Println("Decode request state report fail, ", err)
		return
	}
	defer req.Body.Close()

	str, err := json.Marshal(&state)
	if nil != err {
		g_log.Info.Println("Decode ReqStateReport fail", err)
		return
	}
	g_log.Debug.Println("Update request:", reqType, string(str))
}

func createRequest(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqType := req.FormValue("type")

	body, err := ioutil.ReadAll(req.Body)
	if nil != err {
		g_log.Info.Println("Decode post request body fail, ", err)
		return
	}
	defer req.Body.Close()

	g_log.Debug.Println("Create request:", reqType, string(body))
	addRequest(reqType, string(body))
}

func deleteRequest(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println("form:", req.Form)
}

func requestOpDispatch(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Action:", req.Method)
	switch req.Method {
	case "GET":
		fetchRequest(rw, req)
	case "POST":
		createRequest(rw, req)
	case "PUT":
		updateRequest(rw, req)
	case "DELETE":
		deleteRequest(rw, req)
	default:
		g_log.Info.Println("Unknown request Method: ", req.Method)
	}

}

func clean(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqType := req.FormValue("type")

	cleanRequestStateTable(reqType)
	cleanRequestWaitingQueue(reqType)
	cleanRequestTable(reqType)
}

func setupService() {
	g_log.Info.Println("Listening Port 1234 ...")
	http.HandleFunc("/request", requestOpDispatch)
	http.HandleFunc("/clean", clean)
	http.ListenAndServe(":1234", nil)
}
