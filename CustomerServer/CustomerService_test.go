package main

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCustom(t *testing.T) {

	// Input: ID, Output: expected complete result from MYSQL server.
	testcases := []struct {
		input  int
		output []customer
	}{
		{1, []customer{
			{1, "Sumit", "28/09/1997", address{1, "AKJ", "HSR", "U.P.", 1}},
			{1, "Sumit", "28/09/1997", address{3, "Wallstreet", "NY", "USA.", 1}},
		}},
		{2, []customer{{2, "Sammer", "28/09/1997", address{2, "BKJ", "BTM", "W.B.", 2}}}},

		{0, []customer{

			{1, "Sumit", "28/09/1997", address{1, "AKJ", "HSR", "U.P.", 1}},
			{1, "Sumit", "28/09/1997", address{3, "Wallstreet", "NY", "USA.", 1}},
			{2, "Sammer", "28/09/1997", address{2, "BKJ", "BTM", "W.B.", 2}},
		}},
	}
	// Establishing connection with the database.
	db, err := sql.Open("mysql", "sumit:1234@/customer_service")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
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
