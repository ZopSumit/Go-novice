package main

var u = customer {
	ID  : 1
	Name : "Sumit Saha"
	DOB  : "11-01-1998"
}
import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestSQLQueries(t *testing.T) {
	db, mock := NewMock()
	query := "SELECT * FROM customer where id = \\?"
	rows := sqlmock.NewRows([]string{"ID", "Name", "DOB"}).
		AddRow(u.ID, u.Name, u.DOB)

		mock.ExpectQuery(query).WithArgs(u.ID).WillReturnRows(rows)

		user, err := repo.FindByID(u.ID)
		assert.NotNil(t, user)
		assert.NoError(t, err)
}

