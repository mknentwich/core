package database

import "testing"

func TestInitializeDb(t *testing.T){
	err := InitializeDb()
	if err != nil{
		t.Errorf("Error on creating the database: %s", err.Error())
	}
}
