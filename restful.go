package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReqStateReport struct {
	RequestId string
	WorkerId  string
	State     REQUEST_STATE_TYPE
}

func getRequest(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println("form:", req.Form)
	fmt.Println("reqType:", req.FormValue("type"))
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
		return
	}
	g_log.Debug.Println(reqType, string(str))
}

func createRequest(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqType := req.FormValue("type")

	buff := make([]byte, req.ContentLength)
	body, err := req.Body.Read(buff)
	if nil != err {
		g_log.Info.Println("Decode put request body fail, ", err)
		return
	}
	g_log.Debug.Println(reqType, body)
}

func deleteRequest(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println("form:", req.Form)
}

func requestOpDispatch(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Action:", req.Method)
	switch req.Method {
	case "GET":
		getRequest(rw, req)
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

func setupService() {
	g_log.Info.Println("Listening Port 1234 ...")
	http.HandleFunc("/request", requestOpDispatch)
	http.ListenAndServe(":1234", nil)
}
