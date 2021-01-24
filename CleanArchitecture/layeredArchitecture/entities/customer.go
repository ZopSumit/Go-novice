package entities

type Customer struct {
	ID   int     `json:id`
	Name string  `json:name`
	DOB  string  `json:dob`
	Addr Address `json:addr`
}

type Address struct {
	ID         int    `json:id`
	StreetName string `json:streetName`
	City       string `json:city`
	State      string `json:state`
	CustomerID int    `json:cutomerID`
}
