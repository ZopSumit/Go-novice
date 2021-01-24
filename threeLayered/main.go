package main

import (
	"database/sql"
	"log"
	"net/http"

	"layered/architecture/delivery"
	"layered/architecture/service"
	"layered/architecture/store"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	var db, dbErr = sql.Open("mysql", "sumit:1234@/Cust_Service")
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr.Error()) // proper error handling instead of panic in your app
	}
	r := mux.NewRouter()
	datastore := store.New(db)
	service := service.New(datastore)
	handler := delivery.New(service)

	r.HandleFunc("/customer", handler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/customer", handler.GetByName).Methods(http.MethodGet).Queries("?name","{name:[a-z A-Z 0-9]+}")
	r.HandleFunc("/customer/{id:[0-9]+}", handler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/customer", handler.PostCustomer).Methods(http.MethodPost)
	r.HandleFunc("/customer/{id:[0-9]+}", handler.PutCustomer).Methods(http.MethodPut)
	r.HandleFunc("/customer/{id:[0-9]+}", handler.DeleteCustomer).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8082", r))
}
