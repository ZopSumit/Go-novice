package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	//"golang.org/x/tools/go/analysis"
)

type customer struct {
	ID   int
	Name string
	DOB  string
	Addr address
}
type address struct {
	ID         int
	StreetName string
	City       string
	State      string
	CustomerID int
}

func getCustData(db *sql.DB, c customer) []customer {
	var custValues []interface{}

	query := `INSERT INTO Customer VALUES (?, ?, ?);`
	custValues = append(custValues, c.ID)
	custValues = append(custValues, c.Name)
	custValues = append(custValues, c.DOB)

	_, err := db.Query(query, custValues...)
	if err != nil {
		//panic(err.Error())
		return nil
	}

	// fmt.Println("yes")

	query = `INSERT INTO Address VALUES (?, ?, ?, ?, ?);`
	var addValues []interface{}
	addValues = append(addValues, c.Addr.ID)
	addValues = append(addValues, c.Addr.StreetName)
	addValues = append(addValues, c.Addr.City)
	addValues = append(addValues, c.Addr.State)
	addValues = append(addValues, c.Addr.CustomerID)
	_, err = db.Query(query, addValues...)
	if err != nil {
		//panic(err.Error())
		return nil
	}

	c1 := afterJoining(db)

	return c1
}

func afterJoining(db *sql.DB) []customer {
	var ans []customer
	var idx []interface{}

	query := `SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.CustomerID ORDER BY Customer.ID, Address.ID;`

	// if id != 0 {

	// 	query = `SELECT * FROM Customer INNER JOIN Address ON Customer.ID = Address.CustomerID and CustomerID = ? ORDER BY Customer.ID, Address.ID;`
	// 	idx = append(idx, id)
	// }

	rows, err := db.Query(query, idx...)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var c customer
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.StreetName, &c.Addr.City, &c.Addr.State, &c.Addr.CustomerID); err != nil {
			log.Fatal(err)
		}
		ans = append(ans, c)
	}
	return ans
}

func main() {

}
