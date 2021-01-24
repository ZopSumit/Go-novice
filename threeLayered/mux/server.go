package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Address struct {
	Id         int    `json:"id"`
	StreetName string `json:"streetName"`
	City       string `json:"city"`
	State      string `json:"state"`
	CusId      int    `json:"cusId"`
}

type Customer struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Dob     string  `json:"dob"`
	Address Address `json:"address"`
}

var db, dbErr = sql.Open("mysql", "root:1118209@/Customer_service")

func GetCustomersData(db *sql.DB, name string) []Customer {
	query := "select * from cust inner join addrs on cust.id=addrs.cus_id order by cust.id, addrs.id"
	var data []interface{}
	//fmt.Println("name is ", name)
	if name != "" {
		query = "select * from cust inner join addrs on cust.id=addrs.cus_id where cust.name=? order by cust.id, addrs.id"
		data = append(data, name)
	}

	rows, err := db.Query(query, data...)

	if err != nil {
		return []Customer{}
	}

	var customer []Customer

	for rows.Next() {
		var c Customer
		err = rows.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId)
		customer = append(customer, c)
	}

	return customer
}

func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	//fmt.Println(params)
	name, ok := params["name"]
	var C []Customer
	if ok && len(name) > 0 {
		C = GetCustomersData(db, params["name"][0])
	} else {
		C = GetCustomersData(db, "")
	}
	json.NewEncoder(w).Encode(C)

}

func GetCustomerData(db *sql.DB, id int) Customer {
	rows, err := db.Query("select * from cust inner join addrs on cust.id=addrs.cus_id and cust.id=? order by cust.id, addrs.id", id)

	if err != nil {
		return Customer{}
	}

	var c Customer

	for rows.Next() {
		rows.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId)
	}

	return c
}

func GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	pathparams := mux.Vars(r)

	var c Customer
	id, err := strconv.Atoi(pathparams["id"])

	//pathParams, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Customer{})
	} else {
		c = GetCustomerData(db, id)
		if c.Id == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		json.NewEncoder(w).Encode(c)
	}
}

func InsertCustomerData(db *sql.DB, obj Customer) Customer {
	rows, err := db.Exec("insert into cust (name, dob) values (?, ?)", obj.Name, obj.Dob)

	if err != nil {
		log.Fatal(err)
		return Customer{}
	}

	id, err1 := rows.LastInsertId()

	if err1 != nil {
		log.Fatal(err1)
		return Customer{}
	} else {
		rows, err = db.Exec("insert into addrs (streetname, city, state, cus_id) values (?, ?, ?, ?)", obj.Address.StreetName, obj.Address.City, obj.Address.State, id)
		if err == nil {
			obj.Id = int(id)
			addressId, _ := rows.LastInsertId()
			obj.Address.Id = int(addressId)
			obj.Address.CusId = int(id)
			return obj
		}
	}

	return Customer{}
}

func DateSubstract(d1 string) int {
	d1_slice := strings.Split(d1, "/")

	newDate := d1_slice[2] + "-" + d1_slice[1] + "-" + d1_slice[0]
	myDate, err := time.Parse("2006-01-02", newDate)

	if err != nil {
		//panic(err)
		return 0
	}

	return int(time.Now().Unix() - myDate.Unix())
}

func PostCustomerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var cust Customer
	err = json.Unmarshal(body, &cust)
	if err != nil {
		//log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Customer{})
	} else {
		if cust.Name == "" || cust.Dob == "" || cust.Address.StreetName == "" || cust.Address.City == "" || cust.Address.State == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Customer{})
		} else if timestamp := DateSubstract(cust.Dob); timestamp/(3600*24*12*30) < 18 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Customer{})
		} else {
			cust = InsertCustomerData(db, cust)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(cust)
		}
	}
}

func UpdateData(db *sql.DB, id int, c Customer) Customer {

	if c.Name != "" {
		_, err := db.Exec("update cust set name=? where id=?", c.Name, id)
		if err != nil {
			log.Fatal(err)
			return Customer{}
		}
	}
	var data []interface{}
	query := "update addrs set "
	if c.Address.State != "" {
		query += "state = ? ,"
		data = append(data, c.Address.State)
	}
	if c.Address.City != "" {
		query += "city = ? ,"
		data = append(data, c.Address.City)
	}
	if c.Address.StreetName != "" {
		query += "streetname = ? ,"
		data = append(data, c.Address.StreetName)
	}
	query = query[:len(query)-1]
	query += "where cus_id = ?"
	data = append(data, id)
	_, err := db.Exec(query, data...)

	if err != nil {
		log.Fatal(err)
	}

	rows, _ := db.Query("select * from cust inner join addrs on cust.id=addrs.cus_id and cust.id=?", id)
	var customer Customer
	for rows.Next() {
		rows.Scan(&customer.Id, &customer.Name, &customer.Dob, &customer.Address.Id, &customer.Address.StreetName, &customer.Address.City, &customer.Address.State, &customer.Address.CusId)
	}
	return customer
}

func PutCustomerHandler(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(pathParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Customer{})
	} else {
		var customer Customer
		bodyData, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(bodyData, &customer)

		if err != nil {
			//log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Customer{})
		}
		if (customer.Id != 0 || customer.Dob != "") || (customer.Name == "" && customer.Address.State == "" && customer.Address.City == "" && customer.Address.StreetName == "") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Customer{})
		} else {
			customer = UpdateData(db, id, customer)
			json.NewEncoder(w).Encode(customer)
		}
	}
}

func DeleteData(db *sql.DB, id int) Customer {
	rows, err := db.Query("select * from cust inner join addrs on addrs.cus_id=cust.id and cust.id=? order by cust.id, addrs.id", id)
	if err != nil {
		log.Fatal(err)
		return Customer{}
	}
	var c Customer
	for rows.Next() {
		rows.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId)
	}
	rows, err = db.Query("delete from cust where id=?", id)
	if err != nil {
		log.Fatal(err)
		return Customer{}
	}
	//fmt.Println("in delete ", c)
	return c
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	//pathParams, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	pathParams := mux.Vars(r)["id"]

	id, err := strconv.Atoi(pathParams)
	if err != nil {
		//log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Customer{})
	} else {
		c := DeleteData(db, id)
		if c.Id == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(c)
		} else {
			w.WriteHeader(http.StatusNoContent)
			//fmt.Println(c)
			json.NewEncoder(w).Encode(c)
		}
	}

}

func main() {
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr.Error()) // proper error handling instead of panic in your app
	}
	r := mux.NewRouter()
	r.HandleFunc("/customer", GetCustomersHandler).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", GetCustomerHandler).Methods(http.MethodGet)
	r.HandleFunc("/customer/", PostCustomerHandler).Methods(http.MethodPost)
	r.HandleFunc("/customer/{id}", PutCustomerHandler).Methods(http.MethodPut)
	r.HandleFunc("/customer/{id}", DeleteCustomerHandler).Methods(http.MethodDelete)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))

}
