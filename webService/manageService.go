// manageService.go
package webService

import (
	"context"
	"net/http"
	"scheduler/auth"
	"scheduler/common"
	"scheduler/model"
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

func setupRequestService(ctx context.Context) {

	server := &http.Server{Addr: _managerServAddr_, Handler: nil}

	http.HandleFunc("/addRequestType", addRequestType)
	http.HandleFunc("/clean", clean)

	server.ListenAndServe()
	g_reqHandler.InfoLog("Listening Port 6667 for manage ...")

	go func() {
		select {
		case <-ctx.Done():
			server.Close()
		}
	}()
}

func SetupWebService(ctx context.Context, cephManager *model.CephManager, a *auth.AuthManager) {

	g_cephManager = cephManager
	g_reqHandler = cephManager.GetRequestHandler()

	go setupRequestService(ctx)
	go setupWorkerService(ctx)
	go setupManageService(ctx)
}
