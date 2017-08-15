// manageService.go
package webService

import (
	"context"
	_ "expvar"
	"fmt"
	"net/http"
	"scheduler/common"
	"sync"
)

func setupManageService(wg *sync.WaitGroup, ctx context.Context) {

	server := &http.Server{Addr: _managerServAddr_, Handler: nil}

	http.HandleFunc("/addRequestType", addRequestType)
	http.HandleFunc("/clean", clean)

	g_reqHandler.InfoLog("Listening for manager...on port ", _managerServAddr_)

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Stop Manager Server.")
			server.Close()
			wg.Done()
		}
	}()

	server.ListenAndServe()
}

func addRequestType(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	reqType := req.FormValue("type")

	err := g_reqHandler.CreateNewRequestType(reqType)
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

	err := g_reqHandler.DeleteRequestType(reqType)
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
