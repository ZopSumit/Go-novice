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

func TestGetCustomers(t *testing.T) {

	t.Parallel()
	testCases := []struct {
		inp string
		out []Customer
	}{
		{"?name=Bittu%20Ray", []Customer{{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}}},
		//{"", []Customer{{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}, {Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}, {Id: 3, Name: "Rupesh", Dob: "05-07-2002", Address: Address{Id: 3, StreetName: "rajabazar", City: "Patna", State: "Bihar", CusId: 3}}}},
		//{"?name=", []Customer{{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}, {Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}, {Id: 3, Name: "Rupesh", Dob: "05-07-2002", Address: Address{Id: 3, StreetName: "rajabazar", City: "Patna", State: "Bihar", CusId: 3}}}},
	}

	for i := range testCases {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/customer"+testCases[i].inp, nil)
		w := httptest.NewRecorder()
		GetCustomersHandler(w, req)
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
		{"1", Customer{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}},
		//{"2", Customer{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}},
		{"4", Customer{}},
	}

	for i := range testCases {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/customer/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": testCases[i].inp})
		w := httptest.NewRecorder()
		GetCustomerHandler(w, req)
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

func TestDeleteCustomer(t *testing.T) {
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

//func TestPutCustomer(t *testing.T) {
//	testCases := []struct {
//		inp []byte
//		id  string
//		out Customer
//	}{
//		{[]byte(`{"name":"Pintu","dob":"05-12-2000","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), "3", Customer{Id: 3, Name: "Pintu", Dob: "05-12-2000", Address: Address{Id: 3, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 3}}},
//		//{[]byte(`{}`), Customer{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}},
//		{[]byte(`{"name":"Pintu","dob":"05-12-2000","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), "10", Customer{}},
//		{[]byte(`{"id":1,"name":"Pintu","dob":"05-12-2000","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), "12", Customer(nil)},
//		{[]byte(`{"id":1,"dob":"15-06-2004","name":"Pintu","dob":"05-12-2000","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), "10", Customer(nil)},
//	}
//
//	for i := range testCases {
//		req := httptest.NewRequest(http.MethodPut, "http://localhost:8080/customer/"+testCases[i].id, bytes.NewBuffer(testCases[i].inp))
//		w := httptest.NewRecorder()
//		PutCustomerHandler(w, req)
//		resp := w.Result()
//		body, err := ioutil.ReadAll(resp.Body)
//		var cust Customer
//		err = json.Unmarshal(body, &cust)
//
//		if err != nil {
//			log.Fatal(err)
//		}
//		if !reflect.DeepEqual(cust, testCases[i].out) {
//			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
//		}
//		if http.StatusOK != resp.StatusCode {
//			t.Errorf("FAILED!! expected statusCode %d got %d\n", http.StatusOK, resp.StatusCode)
//		}
//	}
//
//}
