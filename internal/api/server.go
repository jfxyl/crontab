package api

import (
	"crontab/internal/api/controller"
	"crontab/internal/global"
	"fmt"
	"net/http"
)

func InitHttpServer() (err error) {
	var (
		mux           *http.ServeMux
		staticHandler http.Handler
		httpServer    http.Server
		jobController *controller.JobController
	)

	jobController = new(controller.JobController)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", jobController.Save)
	mux.HandleFunc("/job/delete", jobController.Delete)
	mux.HandleFunc("/job/list", jobController.List)
	mux.HandleFunc("/job/executings", jobController.Executings)
	mux.HandleFunc("/job/kill", jobController.Kill)
	mux.HandleFunc("/job/log", jobController.Logs)
	mux.HandleFunc("/worker/list", jobController.WorkList)

	staticHandler = http.FileServer(http.Dir("./web"))
	mux.Handle("/", staticHandler)

	httpServer = http.Server{
		Addr:    fmt.Sprintf(":%d", global.Config.Port),
		Handler: mux,
	}
	err = httpServer.ListenAndServe()
	return
}
