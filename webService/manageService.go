// manageService.go
package webService

import (
	"net/http"
	"scheduler/auth"
	"scheduler/common"
	"scheduler/model"
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

func setupRequestService() {
	g_reqHandler.InfoLog("Listening Port 6667 for manage ...")

	http.HandleFunc("/addRequestType", addRequestType)
	http.HandleFunc("/clean", clean)
	http.ListenAndServe(":6667", nil)
}

func SetupWebService(cephManager *model.CephManager, a *auth.AuthManager) {

	g_cephManager = cephManager
	g_reqHandler = cephManager.GetRequestHandler()

	go setupRequestService()
	go setupWorkerService()
	go setupManageService()
}
