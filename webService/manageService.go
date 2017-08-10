// manageService.go
package webService

import (
	"context"
	_ "expvar"
	"fmt"
	"net/http"
	"scheduler/auth"
	"scheduler/common"
	"scheduler/model"
	"sync"
)

const (
	_requestServAddr_ = ":6666"
	_managerServAddr_ = ":6667"
	_taskServAddr_    = ":6668"
)

var g_auth *auth.AuthManager

var g_cephManager *model.CephManager
var g_reqHandler *model.RequestHandler

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

func setupRequestService(wg *sync.WaitGroup, ctx context.Context) {

	server := &http.Server{Addr: _managerServAddr_, Handler: nil}

	http.HandleFunc("/addRequestType", addRequestType)
	http.HandleFunc("/clean", clean)

	g_reqHandler.InfoLog("Listening Port 6667 for manage ...")

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

func SetupWebService(wg *sync.WaitGroup, ctx context.Context, cephManager *model.CephManager, a *auth.AuthManager) {

	g_cephManager = cephManager
	g_reqHandler = cephManager.GetRequestHandler()

	wg.Add(3)

	c, _ := context.WithCancel(ctx)

	go setupWorkerService(wg, c)
	go setupManageService(wg, c)
	go setupRequestService(wg, c)
}
