package main

import (
	"fmt"
	"net/http"
)

type ReqStateReport struct {
	RequestId string
	WorkerId  string
	State     REQUEST_STATE_TYPE
}

func handleRequest(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println("Action:", req.Method)
	fmt.Println("form:", req.Form)
}

func setupService() {
	http.HandleFunc("/request", handleRequest)
	http.ListenAndServe(":1234", nil)
}
