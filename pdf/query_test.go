package pdf

import (
	"errors"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/rest"
	"net/http"
	"os"
	"reflect"
	"testing"
)

//Starts database service at begin of testing
func TestMain(m *testing.M) {
	go func() {
		context.InitializeCustomConfig(map[string]context.Serve{
			"db": database.Serve,
		},
			&context.Configuration{
				Host:       "0.0.0.0:9400",
				SQLiteFile: ":memory:",
			})
	}()
	for err := errors.New(""); err != nil; _, err = http.Get("http://127.0.0.1:9400/") {
	}
	os.Exit(m.Run())
}

//Tests, if the result from the database has the same values as the local object
func TestPDFOrderResult(t *testing.T) {
	insertTestData()
	result, err := QueryOrderFromIdForPDF(1)
	if err != nil {
		t.Error(err.Error())
	}
	//TODO: get results from database
	expectedResult := OrderResultPDF{
		City:           "Hürth",
		PostCode:       "50354",
		State:          "Deutschland",
		Street:         "Kalscheurener Straße",
		StreetNumber:   "89",
		ID:             1,
		Company:        "Millionen Show",
		Date:           1568024628000,
		FirstName:      "Günter",
		LastName:       "Jauch",
		Salutation:     "Herr",
		ScoreAmount:    3,
		Title:          "Eine letzte Runde (Blasorchesterfassung)",
		Price:          39.9,
		ReferenceCount: 2,
		BillingDate:    20190930,
	}
	if reflect.DeepEqual(result, expectedResult) == false {
		t.Errorf("%s", "Object from the database and local Object aren't the same!")
	}
}

//Inserts TestData to the Database
func insertTestData() {
	address := database.Address{
		City:         "Hürth",
		PostCode:     "50354",
		State:        "Deutschland",
		Street:       "Kalscheurener Straße",
		StreetNumber: "89",
	}
	rest.SaveAddress(address)
	category := database.Category{
		Name: "Polka",
	}
	rest.SaveCategory(category)
	score := database.Score{
		CategoryID: 1,
		Difficulty: 3,
		Price:      39.9,
		Title:      "Eine letzte Runde (Blasorchesterfassung)",
	}
	rest.SaveScore(score)
	order := database.Order{
		BillingAddressID:  1,
		Company:           "Millionen Show",
		Date:              1568024628000,
		DeliveryAddressID: 1,
		Email:             "jauch@werwirdswohl.de",
		FirstName:         "Günter",
		LastName:          "Jauch",
		Payed:             false,
		Salutation:        "Herr",
		ScoreID:           1,
		ScoreAmount:       3,
		Telephone:         "",
		ReferenceCount:    2,
		BillingDate:       20190930,
	}
	rest.SaveOrder(order)
}

//Inserts TestData to the Database
func insertTestData2() {
	address := database.Address{
		City:         "Wien",
		PostCode:     "1050",
		State:        "Österreich",
		Street:       "Spengergasse",
		StreetNumber: "20",
	}
	rest.SaveAddress(address)
	category := database.Category{
		Name: "Polka",
	}
	rest.SaveCategory(category)
	score := database.Score{
		CategoryID: 1,
		Difficulty: 1,
		Price:      24.37,
		Title:      "Über den (Netzwerk)Brücken (Wiener Linien Fassung)",
	}
	rest.SaveScore(score)
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
	rest.SaveOrder(order)
}
