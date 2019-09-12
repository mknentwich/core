package rest

import (
	"fmt"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
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

func TestPDFOrderResult(t *testing.T) {
	insertTestData()
	result := QueryOrderFromIdForPDF(1)
	expectedResult := OrderResultPDF{
		City:           "Hürth",
		PostCode:       "50354",
		State:          "DEUTSCHLAND",
		Street:         "Kalscheurener Straße",
		StreetNumber:   "89",
		ID:             1,
		Company:        "Millionen",
		Date:           1568024628000,
		FirstName:      "Günter",
		LastName:       "Jauch",
		Salutation:     "Herr",
		ScoreAmount:    1,
		Title:          "Eine letzte Runde (Blasorchesterfassung)",
		Price:          39.9,
		ReferenceCount: 2019091202,
	}
	if reflect.DeepEqual(result, expectedResult) == false {
		t.Errorf("%s", "Object from the database and local Object aren't the same!")
	}
}

//Tests, if the highest referenceCount was found
func TestFindMaxReferenceCount(t *testing.T) {
	insertTestData()
	max := FindMaxReferenceCount()
	if max != 2019091202 {
		t.Errorf("%s", "ReferenceCount should be 2019091202, but was "+fmt.Sprint(max))
	}
	insertTestData2()
	max = FindMaxReferenceCount()
	if max != 2019091203 {
		t.Errorf("%s", "ReferenceCount should be 2019091203, but was "+fmt.Sprint(max))
	}
}

//Inserts TestData to the Database
func insertTestData() {
	address := database.Address{
		City:         "Hürth",
		PostCode:     "50354",
		State:        "DEUTSCHLAND",
		Street:       "Kalscheurener Straße",
		StreetNumber: "89",
	}
	InsertNewAddress(address)
	category := database.Category{
		Name: "Polka",
	}
	InsertNewCategory(category)
	score := database.Score{
		CategoryID: 1,
		Difficulty: 3,
		Price:      39.9,
		Title:      "Eine letzte Runde (Blasorchesterfassung)",
	}
	InsertNewScore(score)
	order := database.Order{
		BillingAddressID:  1,
		Company:           "Millionen",
		Date:              1568024628000,
		DeliveryAddressID: 1,
		Email:             "jauch@werwirdswohl.de",
		FirstName:         "Günter",
		LastName:          "Jauch",
		Payed:             false,
		Salutation:        "Herr",
		ScoreID:           1,
		ScoreAmount:       1,
		Telephone:         "",
		ReferenceCount:    2019091202,
	}
	InsertNewOrder(order)
}

//Inserts TestData to the Database
func insertTestData2() {
	address := database.Address{
		City:         "Wien",
		PostCode:     "1050",
		State:        "ÖSTERREICH",
		Street:       "Spengergasse",
		StreetNumber: "20",
	}
	InsertNewAddress(address)
	category := database.Category{
		Name: "Polka",
	}
	InsertNewCategory(category)
	score := database.Score{
		CategoryID: 1,
		Difficulty: 1,
		Price:      24.37,
		Title:      "Über den (Netzwerk)Brücken (Wiener Linien Fassung)",
	}
	InsertNewScore(score)
	order := database.Order{
		BillingAddressID:  1,
		Company:           "",
		Date:              1568024628000,
		DeliveryAddressID: 1,
		Email:             "hpberger@spengergasse.at",
		FirstName:         "Hans-Peter",
		LastName:          "Berger",
		Payed:             false,
		Salutation:        "Herr",
		ScoreID:           1,
		ScoreAmount:       1,
		Telephone:         "",
		ReferenceCount:    2019091203,
	}
	InsertNewOrder(order)
}
