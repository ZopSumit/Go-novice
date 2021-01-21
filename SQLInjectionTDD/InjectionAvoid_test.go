package main

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSQLQueries(t *testing.T) {

	db, err := sql.Open("mysql", "sumit:1234@/Cust_Service")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	query1 := `DROP TABLE Address;`
	query2 := `DROP TABLE Customer;`
	query3 := `CREATE TABLE Customer (ID INT NOT NULL AUTO_INCREMENT, Name VARCHAR(255), DOB VARCHAR(255), PRIMARY KEY(ID));`
	query4 := `CREATE TABLE Address (Id INT NOT NULL AUTO_INCREMENT,StreetName VARCHAR(255), City VARCHAR(255), State VARCHAR(255),CustomerID INT, PRIMARY KEY(Id), FOREIGN KEY (CustomerID) REFERENCES Customer(Id));`

	_, _ = db.Query(query1)
	_, _ = db.Query(query2)
	_, _ = db.Query(query3)
	_, _ = db.Query(query4)

	testcases := []struct {
		input  customer
		output []customer
	}{

		{customer{1, "Sumit", "28/09/1997", address{1, "AKJ", "HSR", "UP", 1}},
			[]customer{
				{1, "Sumit", "28/09/1997", address{1, "AKJ", "HSR", "UP", 1}},
			}},

		{customer{2, "Sammer", "28/09/1997", address{2, "BKJ", "BTM", "WB", 2}},
			[]customer{
				{1, "Sumit", "28/09/1997", address{1, "AKJ", "HSR", "UP", 1}},
				{2, "Sammer", "28/09/1997", address{2, "BKJ", "BTM", "WB", 2}},
			}},
		{
			customer{3, "Bittu", "28/09/1997", address{3, "Wallstreet", "NY", "USA", 1}},
			[]customer{

				{1, "Sumit", "28/09/1997", address{1, "AKJ", "HSR", "UP", 1}},
				{1, "Sumit", "28/09/1997", address{3, "Wallstreet", "NY", "USA", 1}},
				{2, "Sammer", "28/09/1997", address{2, "BKJ", "BTM", "WB", 2}},
			},
		},
	}

	for idx := range testcases {
		actual := getCustData(db, testcases[idx].input)
		if cmp.Equal(actual, testcases[idx].output) == false {
			t.Error("Failed")
			t.Logf("Expected: %v, \n Got: %v", (testcases[idx].output), actual)
			fmt.Println(actual)
		}
	}
}
