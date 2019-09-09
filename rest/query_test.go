package rest

import (
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/pdf"
	"os"
	"reflect"
	"testing"
)

//Starts database service at begin of testing
func TestMain(m *testing.M) {
	context.InitializeCustomConfig(map[string]context.Serve{
		"/db": database.Serve,
	},
		&context.Configuration{
			Host:       "0.0.0.0:9400",
			SQLiteFile: ":memory:",
		})
	os.Exit(m.Run())
}

//Tests inserts of all tables and compares the result of the loaded order with the local one
func TestPDFOrderResult(t *testing.T) {
	address := database.Address{
		City:         "Hürth",
		PostCode:     "50354",
		State:        "DEUTSCHLAND",
		Street:       "Kalscheurener Straße",
		StreetNumber: "89",
	}
	err := InsertNewAddress(address)
	if err != nil {
		t.Errorf("Error on inserting at the addresses table: %s", err.Error())
	}
	category := database.Category{
		Name: "Polka",
	}
	err = InsertNewCategory(category)
	if err != nil {
		t.Errorf("Error on inserting at the categories table: %s", err.Error())
	}
	score := database.Score{
		CategoryID: 1,
		Difficulty: 3,
		Price:      39.9,
		Title:      "Eine letzte Runde (Blasorchesterfassung)",
	}
	err = InsertNewScore(score)
	if err != nil {
		t.Errorf("Error on inserting at the scores table: %s", err.Error())
	}
	order := database.Order{
		BillingAddressID:  1,
		Company:           "Millionen",
		Date:              1568024628000,
		DeliveryAddressID: 1,
		Email:             "jauch@werwirdswohl.de",
		FirstName:         "Günter",
		LastName:          "Jauch",
		Payed:             false,
		ReferenceCount:    0,
		Salutation:        "Herr",
		ScoreID:           1,
		ScoreAmount:       1,
		Telephone:         "",
	}
	err = InsertNewOrder(order)
	if err != nil {
		t.Errorf("Error on inserting at the orders table: %s", err.Error())
	}
	result := QueryOrderFromIdForPDF(1)
	expectedResult := pdf.OrderResultPDF{
		City:         "Hürth",
		PostCode:     "50354",
		State:        "DEUTSCHLAND",
		Street:       "Kalscheurener Straße",
		StreetNumber: "89",
		ID:           1,
		Company:      "Millionen",
		Date:         1568024628000,
		FirstName:    "Günter",
		LastName:     "Jauch",
		Salutation:   "Herr",
		ScoreAmount:  1,
		Title:        "Eine letzte Runde (Blasorchesterfassung)",
		Price:        39.9,
	}
	if reflect.DeepEqual(result, expectedResult) == false {
		t.Errorf("%s", "Object from the database and local Object aren't the same!")
	}
}
