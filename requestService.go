package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
	response := RequestArray{Num: reqNum, RequestList: reqArr}
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

	body, err := ioutil.ReadAll(req.Body)
	if nil != err {
		g_log.Info.Println("Decode post request body fail, ", err)
		return
	}
	defer req.Body.Close()

	g_log.Debug.Println("Create request:", reqType, string(body))
	addRequest(reqType, string(body))
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

func clean(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqType := req.FormValue("type")

	cleanRequestStateTable(reqType)
	cleanRequestWaitingQueue(reqType)
	cleanRequestTable(reqType)
}

func setupRequestService() {
	g_log.Info.Println("Listening Port 1234 for request...")
	http.HandleFunc("/request", requestOpDispatch)
	http.HandleFunc("/clean", clean)
	http.ListenAndServe(":1234", nil)
}
