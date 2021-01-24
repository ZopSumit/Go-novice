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
	c.service.Get(w, id)

}

func (c CustomerHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name, ok := params["name"]

	if !ok {
		c.service.GetByName(w, "")
		return
	}
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
