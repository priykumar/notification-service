package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/priykumar/notification-service/datastore"
	"github.com/priykumar/notification-service/handler"
	"github.com/priykumar/notification-service/service"
)

func main() {

	db := datastore.InitialiseDB()
	tSvc := service.NewTemplateService(db)
	nSvc := service.NewNotificationService(db)
	tHandler := handler.NewTemplateHandler(tSvc)
	nHandler := handler.NewNotificationHandler(nSvc)

	go service.MonitorAndPop(nSvc)

	r := mux.NewRouter()
	r.HandleFunc("/producer/template", tHandler.CreateTemplate).Methods("POST")
	r.HandleFunc("/producer/notify", nHandler.CreateNotification).Methods("POST")

	http.ListenAndServe(":8080", r)
}
