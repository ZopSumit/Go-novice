package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	//"golang.org/x/tools/go/analysis"

	"log"
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

func getCustData(db *sql.DB, id int) []customer {
	var ans []customer
	var idx []interface{}

	query := `SELECT * from Customer INNER JOIN Address ON Customer.id = Address.Customer_id ORDER BY Customer.id, Address.id`

	if id != 0 {

		query = `SELECT * FROM Customer INNER JOIN Address ON Customer.id = Address.Customer_id and Customer.id= ? ORDER BY Customer.id, Address.id`
		idx = append(idx, id)
	}

	rows, err := db.Query(query, idx...)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var c customer
		if err := rows.Scan(&c.ID, &c.Name, &c.DOB, &c.Addr.ID, &c.Addr.Streetname, &c.Addr.City, &c.Addr.State, &c.Addr.Customerid); err != nil {
			log.Fatal(err)
		}
		ans = append(ans, c)
	}
	return ans

}
