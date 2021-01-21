package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetCustByName(t *testing.T) {

	t.Parallel()
	testCases := []struct {
		inp string
		out []Customer
	}{
		{"?name=ZopSumit", []Customer{
			{1, "ZopSumit", "28/09/1997", Address{1, "Athpur", "KOlkata", "WestBangal", 1}}},
		},
	}

	for i := range testCases {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8082/customer"+testCases[i].inp, nil)
		w := httptest.NewRecorder()
		getCustByName(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		var cust []Customer
		err = json.Unmarshal(body, &cust)

		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(cust, testCases[i].out) {
			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
		}
		if http.StatusOK != resp.StatusCode {
			t.Errorf("FAILED!! expected statusCode %d got %d\n", http.StatusOK, resp.StatusCode)
		}
	}

}

func TestGetCustomer(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		inp string
		out Customer
	}{
		{"1", Customer{ID: 1, Name: "ZopSumit", DOB: "28/09/1997", Addr: Address{ID: 1, StreetName: "Athpur", City: "KOlkata", State: "WestBangal", CustomerID: 1}}},
		//{"2", Customer{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}},
		{"4", Customer{}},
	}

	for i := range testCases {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/customer/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": testCases[i].inp})
		w := httptest.NewRecorder()
		getCustByID(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		var cust Customer
		err = json.Unmarshal(body, &cust)

		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(cust, testCases[i].out) {
			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
		}
		if http.StatusOK != resp.StatusCode {
			t.Errorf("FAILED!! expected statusCode %d got %d\n", http.StatusOK, resp.StatusCode)
		}
	}

}

func TestPostCustomer(t *testing.T) {
	testCases := []struct {
		inp []byte
		out Customer
	}{
		{[]byte(`{"name":"Pintu","dob":"05-12-2000","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), Customer{Id: 25, Name: "Pintu", Dob: "05-12-2000", Address: Address{Id: 20, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 25}}},
		//{[]byte(`{"name":"Pintu","dob":"05-12-2000","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), Customer{Id: 4, Name: "Pintu", Dob: "05-12-2000", Address: Address{Id: 4, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 4}}},
		//{[]byte(`{"name":"Pintu","dob":"05-12-2006","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), Customer(nil)},
		//{[]byte(`{}`), Customer{}},
	}

	for i := range testCases {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/customer/", bytes.NewBuffer(testCases[i].inp))
		w := httptest.NewRecorder()
		PostCustomerHandler(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		var cust Customer
		err = json.Unmarshal(body, &cust)

		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(cust, testCases[i].out) {
			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
		}

		if http.StatusCreated != resp.StatusCode {
			t.Errorf("FAILED!! expected statusCode %d got %d\n", http.StatusCreated, resp.StatusCode)
		}
	}

}

func TestDeleteByID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		inp string
		out Customer
	}{
		{"3", Customer{Id: 3, Name: "Rupesh", Dob: "05-07-2002", Address: Address{Id: 3, StreetName: "rajabazar", City: "Patna", State: "Bihar", CusId: 3}}},
		//{"2", Customer{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}},
		{"4", Customer{}},
	}

	for i := range testCases {
		req := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/customer/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": testCases[i].inp})
		w := httptest.NewRecorder()
		DeleteCustomerHandler(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		var cust Customer
		err = json.Unmarshal(body, &cust)

		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(cust, testCases[i].out) {
			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
		}

		if http.StatusNoContent != resp.StatusCode {
			t.Errorf("FAILED!! expected statusCode %d got %d\n", http.StatusOK, resp.StatusCode)
		}

	}

}
