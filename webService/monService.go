package webService

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

func setupMonService(wg *sync.WaitGroup, ctx context.Context) {

	server := &http.Server{Addr: _requestServAddr_, Handler: nil}

	http.HandleFunc("/monData", monDataOpDispatch)

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Stop Monitor Server.")
			server.Close()
			wg.Done()
		}
	}()
	g_reqHandler.InfoLog("Listening for request...on port ", _monitorServAddr_)
	server.ListenAndServe()
}

func monDataOpDispatch(rw http.ResponseWriter, req *http.Request) {

	fmt.Println("Action:", req.Method)

	switch req.Method {
	//	case "GET":
	//	case "POST":
	case "PUT":
		insertNewMonData(rw, req)

		//	case "DELETE":

	default:
		s := fmt.Sprintf("It's illegal to %s task!", req.Method)
		g_reqHandler.InfoLog(s)
		rw.Write([]byte(s))
	}
}

func insertNewMonData(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("insert new monitor data")
}
