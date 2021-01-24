package entities

type Address struct {
	ID         int    `json:"id"`
	StreetName string `json:"streetName"`
	City       string `json:"city"`
	State      string `json:"state"`
	CusId      int    `json:"cusId"`
}

type Customer struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	DOB     string  `json:"dob"`
	Address Address `json:"address"`
}
