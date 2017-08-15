package webService

import (
	"context"
	_ "expvar"
	"scheduler/auth"
	"scheduler/model"
	"sync"
)

const (
	_monitorServAddr_ = ":6665"
	_requestServAddr_ = ":6666"
	_managerServAddr_ = ":6667"
	_taskServAddr_    = ":6668"
)

var g_auth *auth.AuthManager

var g_monManager *model.MonitorManager
var g_cephManager *model.CephManager
var g_reqHandler *model.RequestHandler

func SetupWebService(wg *sync.WaitGroup, ctx context.Context, cephManager *model.CephManager, a *auth.AuthManager, m *model.MonitorManager) {

	g_monManager = m

	g_cephManager = cephManager
	g_reqHandler = cephManager.GetRequestHandler()

	wg.Add(3)

	c, _ := context.WithCancel(ctx)

	go setupManageService(wg, c)
	go setupMonService(wg, c)
	go setupRequestService(wg, c)
	go setupWorkerService(wg, c)
}
