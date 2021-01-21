package main

import (
	_ "github.com/go-sql-driver/mysql"
	//"golang.org/x/tools/go/analysis"
)

type customer struct {
	ID   string
	Name string
	DOB  string
}

func FindByID(id string)
