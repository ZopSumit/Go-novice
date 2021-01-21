package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type customer struct {
	ID   int
	Name string
	DOB  string
	Addr address
}
type address struct {
	ID         int
	Streetname string
	City       string
	State      string
	Customerid int
}

func handler(w http.ResponseWriter, r *http.Request) {
	var ans []customer

	db, err := sql.Open("mysql", "sumit:1234@/customer_service?multiStatements=true")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	a := string(strings.Split(r.URL.Path, "/")[1])

	// if len(a) == 0 {
	// 	out, err := db.Query(fmt.Sprintf("SELECT * FROM Customer INNER JOIN Address ON Customer.id = Address.Customer_id ORDER BY Customer.id, Address.id;"))
	// 	if err != nil {
	// 		panic(err.Error()) // proper error handling instead of panic in your app
	// 	}
	// 	for out.Next() {

	// 		var c customer
	// 		if err := out.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.Streetname, &c.Addr.City, &c.Addr.State, &c.Addr.Customerid); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		//fmt.Printf("Print All the %v\n", c)
	// 		ans = append(ans, c)

	// 	}

	// } else {
	// 	idx, err := strconv.Atoi(a)
	// 	out, err := db.Query(fmt.Sprintf("SELECT * FROM Customer INNER JOIN Address ON Customer.id = Address.Customer_id WHERE Customer.id=%v ORDER BY Customer.id, Address.id;", idx))
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	for out.Next() {
	// 		var c customer
	// 		if err := out.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.Streetname, &c.Addr.City, &c.Addr.State, &c.Addr.Customerid); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		//fmt.Println(c)
	// 		ans = append(ans, c)

	// 	}

	// }

	query := "SELECT * FROM Customer INNER JOIN Address ON Customer.id = Address.Customer_id ORDER BY Customer.id, Address.id;" + a

	out, err := db.Query(query)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for out.Next() {

		var c customer
		if err := out.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.Streetname, &c.Addr.City, &c.Addr.State, &c.Addr.Customerid); err != nil {
			log.Fatal(err)
		}
		ans = append(ans, c)

	}

	err = json.NewEncoder(w).Encode(ans)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Main spawns the server.
func main() {

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
