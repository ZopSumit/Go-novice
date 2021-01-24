package store

import "C"
import (
	"database/sql"
	"layered/architecture/entities"
)

type CustomerStore struct {
	db *sql.DB
}

func New(db *sql.DB) Customer {
	return CustomerStore{db: db}
}

func (c CustomerStore) GetByID(id int) (entities.Customer, error) {
	rows, err := c.db.Query("select * from Customer inner join Address on Customer.ID=Address.Customer_id and Customer.ID=? order by Customer.ID, Address.ID", id)
	if err != nil {
		return entities.Customer{}, err
	}

	var cust entities.Customer

	for rows.Next() {
		rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}

	return cust, nil
}
func (c CustomerStore) GetAll() ([]entities.Customer, error) {
	query := "select * from Customer inner join Address on Customer.ID=Address.Customer_id order by Customer.ID, Address.ID"
	rows, err := c.db.Query(query)

	if err != nil {
		return []entities.Customer(nil), err
	}
	var customer []entities.Customer

	for rows.Next() {
		var cust entities.Customer
		err = rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
		customer = append(customer, cust)
	}

	return customer, nil
}
func (c CustomerStore) GetByName(name string) ([]entities.Customer, error) {
	query := "select * from Customer inner join Address on Customer.ID=Address.Customer_id order by Customer.ID, Address.ID"
	var data []interface{}
	if name != "" {
		query = "select * from Customer inner join Address on Customer.ID=Address.Customer_id where Customer.Name=? order by Customer.ID, Address.ID"
		data = append(data, name)
	}

	rows, err := c.db.Query(query, data...)

	if err != nil {
		return []entities.Customer(nil), err
	}
	var customer []entities.Customer

	for rows.Next() {
		var cust entities.Customer
		err = rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
		customer = append(customer, cust)
	}

	return customer, nil
}

func (c CustomerStore) Create(customer entities.Customer) (entities.Customer, error) {

	rows, err := c.db.Exec("insert into Customer (Name, DOB) values (?, ?)", customer.Name, customer.DOB)
	if err != nil {
		return entities.Customer{}, err
	}

	id, err1 := rows.LastInsertId()

	if err1 != nil {
		return entities.Customer{}, err
	} else {
		rows, err = c.db.Exec("insert into Address (StreetName, City, State, Customer_id) values (?, ?, ?, ?)", customer.Address.StreetName, customer.Address.City, customer.Address.State, id)
		if err == nil {
			customer.ID = int(id)
			addressId, _ := rows.LastInsertId()
			customer.Address.ID = int(addressId)
			customer.Address.CusId = int(id)
			return customer, nil
		}
		return entities.Customer{}, err
	}
}

func (c CustomerStore) Update(id int, customer entities.Customer) (entities.Customer, error) {
	if customer.Name != "" {
		_, err := c.db.Exec("update Customer set Name=? where ID=?", customer.Name, id)
		if err != nil {
			return entities.Customer{}, nil
		}
	}
	var data []interface{}
	query := "update Address set "
	if customer.Address.State != "" {
		query += "State = ? ,"
		data = append(data, customer.Address.State)
	}
	if customer.Address.City != "" {
		query += "City = ? ,"
		data = append(data, customer.Address.City)
	}
	if customer.Address.StreetName != "" {
		query += "StreetName = ? ,"
		data = append(data, customer.Address.StreetName)
	}
	query = query[:len(query)-1]
	query += "where Customer_id = ?"
	data = append(data, id)
	_, err := c.db.Exec(query, data...)

	if err != nil {
		return entities.Customer{}, err
	}

	rows, _ := c.db.Query("select * from Customer inner join Address on Customer.ID=Address.Customer_id and Customer.ID=?", id)
	var cust entities.Customer
	for rows.Next() {
		rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}
	return cust, nil
}

func (c CustomerStore) Delete(id int) (entities.Customer, error) {
	rows, err := c.db.Query("select * from Customer inner join Address on Address.Customer_id=Customer.ID and Customer.ID=? order by Customer.ID, Address.ID", id)
	if err != nil {
		return entities.Customer{}, err
	}
	var cust entities.Customer
	for rows.Next() {
		rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}
	_, err = c.db.Exec("delete from Customer where ID=?", id)
	if err != nil {
		return entities.Customer{}, err
	}
	return cust, nil
}
