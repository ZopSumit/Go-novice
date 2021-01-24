package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	age "github.com/bearbin/go-age"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type customer struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	DOB  string  `json:"dob"`
	Addr address `json:"addr"`
}

type address struct {
	ID         int    `json:"id"`
	StreetName string `json:"street_name"`
	City       string `json:"city"`
	State      string `json:"state"`
	CusID      int    `json:"cus_id"`
}

// Router global, so that handlers can access the routers
var router = mux.NewRouter()
var db *sql.DB
var err error

func connectDatabase() {
	db, err = sql.Open("mysql", "root:password@/customer_service")
	fmt.Println("Database connected.")
	if err != nil {
		panic(err)
	}
}

func getCustomerAll(w http.ResponseWriter, r *http.Request) {
	// Return all customers.
	var res []customer
	//w.Header().Set("Content-Type", "application/json")
	query := `SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.Customer_id ORDER BY Customer.ID, Address.ID`
	rows, err := db.Query(query)
	if err != nil {
		// Learn to handle panic.
		panic(err.Error())
	}

	// Iterate through all customers.
	for rows.Next() {
		var c customer
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CusID); err != nil {
			log.Fatal(err)
		}
		res = append(res, c)
	}

	_ = json.NewEncoder(w).Encode(res)
}

func getCustomerByID(w http.ResponseWriter, r *http.Request) {
	// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
	//w.Header().Set("Content-Type", "application/json")
	// Take all variables in the multiplexer as params.
	// That allows us to pass that variable in JSON as ? in query param of SQL.
	params := mux.Vars(r)

	row := db.QueryRow("SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.Customer_id and Customer.ID = ? ORDER BY Customer.ID, Address.ID;", params["id"])

	//for rows.Next() {
	var c customer
	if err := row.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CusID); err != nil {
		log.Fatal(err)
	}
	//	res = append(res, c)
	//}

	// JSON data, combines marshal and writer.
	_ = json.NewEncoder(w).Encode(c)
}

func getCustomerByName(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	// Take all variables in the multiplexer as params.
	// That allows us to pass that variable in JSON as ? in query param of SQL.
	params := mux.Vars(r)
	rows, err := db.Query("SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.Customer_id and Customer.Name = ? ORDER BY Customer.ID, Address.ID;", params["name"])
	if err != nil {
		panic(err.Error())
	}
	var res []customer // Changing tc
	for rows.Next() {
		var c customer
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CusID); err != nil {
			log.Fatal(err)
		}

		res = append(res, c)
	}
	_ = json.NewEncoder(w).Encode(res)

}

// func deleteCustomer, which deletes /customer/{id}. in the database belonging to an ID.
// In the function, just return string "SUCCESS" in the implementation.
// We dont want to delete the data.
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Setting up content tye for the response to JSON.
	w.Header().Set("Content-Type", "application/json")
	// Query parameters are route variables for the current request.
	params := mux.Vars(r)

	// Preparing a query to be executed. (May execute multiple statements here as well)
	// Get the data
	getCustomerByID(w, r)
	// Delete not working on DB.
	stmt, err := db.Prepare("DELETE FROM Customer WHERE ID = ?;")
	if err != nil {
		panic(err.Error())
	}
	_, _ = stmt.Exec(params["id"])
	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
	//io.WriteString(os.Stdout, fmt.Sprintf("Post with ID = %s was deleted", params["id"]))
	// Internally converts given string to []byte/ JSON.
	_, _ = io.WriteString(w, "SUCCESS")
}

//--------------------------------------------------------------------------------------------------

// AgeAt gets the age of an entity at a certain time.
func AgeAt(birthDate time.Time, now time.Time) int {
	// Get the year number change since the player's birth.
	years := now.Year() - birthDate.Year()

	// If the date is before the date of birth, then not that many years have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years--
	}

	return years
}

// Age is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now())
}

// Gets the adjusted date of birth to work around leap year differences.
func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

// Works out if a time.Time is in a leap year.
func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}

// func calculateAge(age string) int {
// 	// https://golangbyexample.com/get-age-given-dob-go/
// 	temp := strings.Split(age, "-")
// 	year, _ := strconv.Atoi(temp[2])
// 	month, _ := strconv.Atoi(temp[1])
// 	day, _ := strconv.Atoi(temp[0])
// 	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
// 	return Age(dob)
// }

// ---------------------------------------------------------------------------------------------------

// // POST should add a customer, if already exists.
// /*func postCustomer(w http.ResponseWriter, r *http.Request) {
// 	// reads bytes from input output.
// 	body, _ := ioutil.ReadAll(r.Body)
// 	// Using unmarshal, convert JSON/[]byte to struct as required.
// 	var c customer
// 	// During unmarshal, only named values will be filled.
// 	_ = json.Unmarshal(body, &c)
// 	/*
// 		Calculate age from current date in GoLang.
// 		If in DB, age is not given, calculate current age using DOB.
// 	*/
// 	age := calculateAge(c.DOB)
// 	if age <= 18 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		_, _ = io.WriteString(w, "NOT ALLOWED")
// 		return
// 	}
// 	if len(c.Name) == 0 || len(c.DOB) == 0 || len(c.Addr.StreetName) == 0 || len(c.Addr.City) == 0 || len(c.Addr.State) == 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	// Create a new record
// 	var cust, addr []interface{}
// 	cust = append(cust, c.Name)
// 	cust = append(cust, c.DOB)
// 	query1 := `INSERT INTO Customer(Name, DOB) VALUES(?, ?);`
// 	// https://medium.com/@alok.sinha.nov/query-vs-exec-vs-prepare-in-golang-e7c49212c36c
// 	rows1, err1 := db.Exec(query1, cust...)
// 	if err1 != nil {
// 		panic(err1.Error())
// 	}
// 	// Get the ID of the last inserted value, because we will need its id to get the value later from DB.
// 	id1, _ := rows1.LastInsertId()
// 	c.ID = int(id1)

// 	addr = append(addr, c.Addr.StreetName)
// 	addr = append(addr, c.Addr.City)
// 	addr = append(addr, c.Addr.State)
// 	addr = append(addr, id1)
// 	query2 := `INSERT INTO Address(StreetName, City, State, Customer_id) VALUES(?, ?, ?, ?);`
// 	rows2, err2 := db.Exec(query2, addr...)
// 	if err2 != nil {
// 		panic(err2.Error())
// 	}
// 	id2, _ := rows2.LastInsertId()

// 	// Update c.ID and c.Addr.ID
// 	c.Addr.ID = int(id2)

// 	//fmt.Println(c)
// 	//byte, _ := json.Marshal(c)
// 	//w.Write(byte)
// 	_ = json.NewEncoder(w).Encode(c)

// 	// In default case, return empty slice.
// 	// json.NewEncoder(w).Encode(customer{})
// }

func getDOB(year, month, day int) time.Time {

	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dob
}

func postCustData(w http.ResponseWriter, r *http.Request) {

	var c customer
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
	//fmt.Println(age.Age(getAge))
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

		var ans []customer
		c.ID = int(idAddr)
		c.Addr.CusID = int(idAddr)
		for row.Next() {
			if err := row.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CusID); err != nil {
				log.Fatal(err)
			}
			ans = append(ans, c)
		}
		w.WriteHeader(http.StatusCreated)
		Byte, _ := json.Marshal(ans)
		_, err = io.WriteString(w, string(Byte))
		//json.NewEncoder(w).Encode(ans)
	} else {
		w.WriteHeader(http.StatusBadRequest)

	}
}

func putCustomer(w http.ResponseWriter, r *http.Request) {
	var c customer
	_ = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	param := mux.Vars(r)
	id := param["id"]
	fmt.Println(id)
	var data1 []interface{}

	if c.Name != "" {
		query := "update Customer set Name=? where ID=?"
		fmt.Println(query)
		data1 = append(data1, c.Name)
		data1 = append(data1, id)

		_, err := db.Exec(query, data1...)
		fmt.Println("sumit")
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println(c)
	var data2 []interface{}
	query := "UPDATE Address set "
	fmt.Println(c.Addr.City, c.Addr.StreetName, c.Addr.State)
	if len(c.Addr.City) != 0 {
		query = query + "City = ? ,"
		data2 = append(data2, c.Addr.City)
	}
	if len(c.Addr.State) != 0 {
		query = query + "State = ? ,"
		data2 = append(data2, c.Addr.State)
	}
	if len(c.Addr.StreetName) != 0 {
		query = query + "StreetName = ? ,"
		data2 = append(data2, c.Addr.StreetName)
	}
	query = query[:len(query)-1]
	query = query + " where Customer_id =? and Address.ID =?; "
	data2 = append(data2, id)
	data2 = append(data2, c.Addr.ID)
	fmt.Println(query)
	_, err = db.Exec(query, data2...)

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(c)
}

/*
//Update customer values,
func putCustomer(w http.ResponseWriter, r *http.Request) {
	// Connect the db.
	connectDatabase()

	//https://golangbyexample.com/net-http-package-get-query-params-golang/
	var c customer
	_ = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	// Cant update DOB.
	if len(c.DOB) > 0 {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Check for valid DOB.
	params := mux.Vars(r)
	id := params["id"]
	var firstItr []interface{}

	if c.Name != "" {
		query := "update Customer set Name=? where ID=?"
		firstItr = append(firstItr, c.Name)
		firstItr = append(firstItr, id)
		_, err := db.Exec(query, firstItr...)
		if err != nil {
			panic(err.Error())
		}
	}

	var secondItr []interface{}
	query := "UPDATE Address set "
	if c.Addr.City != "" {
		query += "City = ? ,"
		secondItr = append(secondItr, c.Addr.City)
	}
	if c.Addr.State != "" {
		query += "State = ? ,"
		secondItr = append(secondItr, c.Addr.State)
	}
	if c.Addr.StreetName != "" {
		query += "StreetName = ? ,"
		secondItr = append(secondItr, c.Addr.StreetName)
	}
	// Remove last ,
	query = query[:len(query)-1]
	// Add comparison for foreign key.
	query += " where Address.CusID = ? and Address.ID = ?" // Single row.
	secondItr = append(secondItr, id)
	secondItr = append(secondItr, c.Addr.ID)
	_, err = db.Exec(query, secondItr...)

	if err != nil {
		log.Fatal(err)
	}
	//getCustomerByID(w, r)
	json.NewEncoder(w).Encode(c)
}
*/

func main() {
	db, err = sql.Open("mysql", "sumit:1234@/Cust_Service")
	fmt.Println("Database connected.")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Retrieve
	router.HandleFunc("/customer", getCustomerAll).Methods(http.MethodGet)
	router.HandleFunc("/customer/{id:[0-9]+}", getCustomerByID).Methods(http.MethodGet)
	router.HandleFunc("/customer/{name:[a-zA-Z ]+}", getCustomerByName).Methods(http.MethodGet)
	// Delete
	router.HandleFunc("/customer/{id:[0-9]+}", deleteCustomer).Methods(http.MethodDelete)
	// Create
	router.HandleFunc("/customer", postCustData).Methods(http.MethodPost)
	// Update
	router.HandleFunc("/customer/{id:[0-9]+}", putCustomer).Methods(http.MethodPut)

	if err := http.ListenAndServe(":8082", router); err != nil {
		// Handle error properly in your app.
		err1 := errors.New("problem spawning port")
		log.Fatal(err1, err)
	}

}
