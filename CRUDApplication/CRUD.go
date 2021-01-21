package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	age "github.com/bearbin/go-age"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Customer struct {
	ID   int    `json:id`
	Name string `json:name`
	DOB  string `json:dob`
	Addr Address
}

type Address struct {
	ID         int    `json:id`
	StreetName string `json:streetname`
	City       string `json:city`
	State      string `json:state`
	CustomerID int    `json:cutomerid`
}

var db *sql.DB
var err error

// METHOD : GET by name or empty

func getCustByName(w http.ResponseWriter, r *http.Request) {
	var ans []Customer
	q := r.URL.Query()
	var data []interface{}
	name, ok := q["name"]
	query := `SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.Customer_id ORDER BY Customer.ID, Address.ID;`

	if ok {
		query = `SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.Customer_id and Customer.Name = ? ORDER BY Customer.ID, Address.ID;`
		data = append(data, name[0])
	}

	rows, err := db.Query(query, data...)

	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CustomerID); err != nil {
			log.Fatal(err)
		}
		ans = append(ans, c)
	}

	json.NewEncoder(w).Encode(ans)

}

// METHOD : GET by ID
func getCustByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var data []interface{}
	var ans []Customer

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	//fmt.Printf("yes after error id is %T, %v ", id, id)

	query := `SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.Customer_id and Customer.ID = ? ORDER BY Customer.ID, Address.ID;`
	data = append(data, id)
	rows, err := db.Query(query, data...)
	if err != nil {
		panic(err.Error())
	}
	var c Customer
	for rows.Next() {
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CustomerID); err != nil {
			log.Fatal(err)
		}
		//fmt.Println("in loop c is ", c)
		ans = append(ans, c)
	}
	//fmt.Println("ans is ", ans)
	if c.Name == "" {
		//fmt.Println("yes in Name")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	//fmt.Println("yes at end")
	json.NewEncoder(w).Encode(ans)
}

func getDOB(year, month, day int) time.Time {

	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dob
}

// METHOD : POST
func postCustData(w http.ResponseWriter, r *http.Request) {

	var c Customer
	err = json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		panic(err.Error())
	}

	if len(c.DOB) == 0 || len(c.Name) == 0 || len(c.Addr.City) == 0 || len(c.Addr.State) == 0 || len(c.Addr.StreetName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	dob := c.DOB
	x := strings.Split(dob, "/")

	y, _ := strconv.Atoi(x[2])
	m, _ := strconv.Atoi(x[1])
	d, _ := strconv.Atoi(x[0])
	getAge := getDOB(y, m, d)
	fmt.Println(age.Age(getAge))
	if age.Age(getAge) >= 18 {

		var custValues []interface{}
		query := `INSERT INTO Customer VALUES (?, ?, ?);`
		custValues = append(custValues, c.ID)
		custValues = append(custValues, c.Name)
		custValues = append(custValues, c.DOB)

		rows, err := db.Exec(query, custValues...)
		if err != nil {
			panic(err.Error())

		}
		idAddr, _ := rows.LastInsertId()
		query = `INSERT INTO Address VALUES (?, ?, ?, ?, ?);`
		var addValues []interface{}
		addValues = append(addValues, c.Addr.ID)
		addValues = append(addValues, c.Addr.StreetName)
		addValues = append(addValues, c.Addr.City)
		addValues = append(addValues, c.Addr.State)
		addValues = append(addValues, idAddr)
		a, _ := db.Exec(query, addValues...)
		id, err := a.LastInsertId()
		if err != nil {
			panic(err.Error())

		}

		query = `SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.Customer_id and Customer.ID = ? and Address.ID=? ORDER BY Customer.ID, Address.ID;`
		var data []interface{}
		data = append(data, int(idAddr))
		data = append(data, int(id))
		row, _ := db.Query(query, data...)

		var ans []Customer
		c.ID = int(idAddr)
		c.Addr.CustomerID = int(idAddr)
		for row.Next() {
			if err := row.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CustomerID); err != nil {
				log.Fatal(err)
			}
			ans = append(ans, c)
		}
		json.NewEncoder(w).Encode(ans)
	} else {
		w.WriteHeader(http.StatusBadRequest)

	}
}

func updateCustData(w http.ResponseWriter, r *http.Request) {
	var c Customer
	body, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(body, &c)
	if err != nil {
		log.Fatal(err)
	}
	param := mux.Vars(r)
	id := param["id"]
	var data1 []interface{}

	if c.Name != "" {
		query := "update Customer set Name=? where ID=?"
		data1 = append(data1, c.Name)
		data1 = append(data1, id)

		_, err := db.Exec(query, data1...)

		if err != nil {
			panic(err.Error())
		}
	}
	var data2 []interface{}
	query := "UPDATE Address set "
	if c.Addr.City != "" {
		query += "City = ? ,"
		data2 = append(data2, c.Addr.City)
	}
	if c.Addr.State != "" {
		query += "State = ? ,"
		data2 = append(data2, c.Addr.State)
	}
	if c.Addr.StreetName != "" {
		query += "StreetName = ? ,"
		data2 = append(data2, c.Addr.StreetName)
	}
	query = query[:len(query)-1]
	query += " where Customer_id =? and ID =?"
	data2 = append(data2, id)
	data2 = append(data2, c.Addr.ID)
	//fmt.Println("query is ", query, " and interface is ", data2)
	_, err = db.Exec(query, data2...)

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(c)

}

// METHOD : DELETE POST
func deleteCustData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	getCustByID(w, r)
	stmt, err := db.Prepare("DELETE FROM Customer WHERE ID = ?;")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Println("yes in nil")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	db, err = sql.Open("mysql", "sumit:1234@/Cust_Service")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/customer", getCustByName).Methods("GET")
	router.HandleFunc("/customer/{id:[0-9]+}", getCustByID).Methods("GET")
	router.HandleFunc("/customer", postCustData).Methods("POST")
	router.HandleFunc("/customer/{id:[0-9]+}", updateCustData).Methods("PUT")
	router.HandleFunc("/customer/{id:[0-9]+}", deleteCustData).Methods("DELETE")
	http.ListenAndServe(":8082", router)
}
