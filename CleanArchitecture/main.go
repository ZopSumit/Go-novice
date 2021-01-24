package main

import (
	delivery "/home/epsilon/Desktop/Layered-Architecture/CleanArchitecture/layeredArchitecture/delivery"
	service "/home/epsilon/Desktop/Layered-Architecture/CleanArchitecture/layeredArchitecture/service"
	"database/sql"

	datastore "/home/epsilon/Desktop/Layered-Architecture/CleanArchitecture/layeredArchitecture/store"

	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	var db, err = sql.Open("mysql", "root:sumit@1234@/Cust_Service")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	r := mux.NewRouter()
	datastore := datastore.New(db)
	service := service.New(datastore)
	handler := delivery.New(service)

	r.HandleFunc("/customer", handler.GetCustByName).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", handler.GetCustById).Methods(http.MethodGet)
	r.HandleFunc("/customer", handler.postCustData).Methods(http.MethodPost)
	r.HandleFunc("/customer/{id}", handler.updateCustData).Methods(http.MethodPut)
	r.HandleFunc("/customer/{id}", handler.deleteCustData).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8080", r))
}
