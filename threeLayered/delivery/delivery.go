package delivery

import (
	"encoding/json"
	"io/ioutil"
	"layered/architecture/entities"
	"layered/architecture/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	service service.CustomerService
}

func New(customer service.CustomerService) CustomerHandler {
	return CustomerHandler{service: customer}
}

func (c CustomerHandler) GetById(w http.ResponseWriter, r *http.Request) {
	pathparams := mux.Vars(r)

	id, err := strconv.Atoi(pathparams["id"])
	if err != nil {
		_, _ = w.Write([]byte("invalid parameter id"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c.service.GetById(w, id)

}

func (c CustomerHandler) GetAll(w http.ResponseWriter, r *http.Request) {

		c.service.GetAll(w)

}

func (c CustomerHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name, _ := params["name"]

	c.service.GetByName(w, name[0])
}

func (c CustomerHandler) PostCustomer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var cust entities.Customer
	err = json.Unmarshal(body, &cust)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid data format"))
		return
	}
	c.service.CreateCustomer(w, cust)
}

func (c CustomerHandler) PutCustomer(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(pathParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id parameter"))
		return
	}
	var customer entities.Customer
	bodyData, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(bodyData, &customer)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Data format is not correct"))
		return
	}
	c.service.UpdateCustomer(w, id, customer)

}

func (c CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)["id"]

	id, err := strconv.Atoi(pathParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id parameter"))
		return
	}
	c.service.DeleteCustomer(w, id)
}
