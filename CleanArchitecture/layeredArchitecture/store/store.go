package datastore

//import "C"
import (
	"database/sql"
	"layered/architecture/entities"
)

type CustomerStore struct {
	db *sql.DB
}

func New(db *sql.DB) CustomerStore {
	return CustomerStore{db: db}
}

func (c CustomerStore) GetById(id int) (entities.Customer, error) {
	rows, err := c.db.Query("select * from cust inner join addrs on cust.id=addrs.cus_id and cust.id=? order by cust.id, addrs.id", id)
	if err != nil {
		return entities.Customer{}, err
	}

	var cust entities.Customer

	for rows.Next() {
		rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}

	return cust, nil
}

func (c CustomerStore) GetByName(name string) ([]entities.Customer, error) {
	query := "select * from cust inner join addrs on cust.id=addrs.cus_id order by cust.id, addrs.id"
	var data []interface{}
	if name != "" {
		query = "select * from cust inner join addrs on cust.id=addrs.cus_id where cust.name=? order by cust.id, addrs.id"
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

	rows, err := c.db.Exec("insert into cust (name, dob) values (?, ?)", customer.Name, customer.DOB)
	if err != nil {
		return entities.Customer{}, err
	}

	id, err1 := rows.LastInsertId()

	if err1 != nil {
		return entities.Customer{}, err
	} else {
		rows, err = c.db.Exec("insert into addrs (streetname, city, state, cus_id) values (?, ?, ?, ?)", customer.Address.StreetName, customer.Address.City, customer.Address.State, id)
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
