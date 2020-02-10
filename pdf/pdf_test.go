package pdf

import (
	"errors"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/rest"
	"net/http"
	"os"
	"testing"
)

var address database.Address
var order database.Order
var score database.Score

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
func TestPDF(t *testing.T) {
	t.Run("InsertTestData_1", func(t *testing.T) {
		insertTestData()
		result, err := QueryOrderFromIdForPDF(1)
		if err != nil {
			t.Error(err.Error())
		}
		if result.ReferenceCount != 1 {
			t.Errorf("%s", "Wrong ReferenceCount from Trigger")
		}
	})
	t.Run("InsertTestData_2", func(t *testing.T) {
		insertTestData2()
		result, err := QueryOrderFromIdForPDF(2)
		if err != nil {
			t.Error(err.Error())
		}
		if result.ReferenceCount != 2 {
			t.Errorf("%s", "Wrong ReferenceCount from Trigger")
		}
	})
	t.Run("InsertTestData_3", func(t *testing.T) {
		insertTestData3()
		result, err := QueryOrderFromIdForPDF(3)
		if err != nil {
			t.Error(err.Error())
		}
		if result.ReferenceCount != 3 {
			t.Errorf("%s", "Wrong ReferenceCount from Trigger")
		}
	})
	t.Run("PDFCreationFromTestData", func(t *testing.T) {
		f, err := os.OpenFile("example-bill-company.pdf", os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
		write, _, err := writeBill(1)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
			return
		}
		err = write(f)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
		f, err = os.OpenFile("example-bill-customer.pdf", os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
		write, _, err = writeBill(2)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
			return
		}
		err = write(f)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
		f, err = os.OpenFile("example-bill-customer2.pdf", os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
		write, _, err = writeBill(3)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
			return
		}
		err = write(f)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
	})
}

//Inserts TestData to the Database
func insertTestData() {
	address = database.Address{
		City:         "Hürth",
		PostCode:     "50354",
		State:        "Deutschland",
		Street:       "Kalscheurener Straße",
		StreetNumber: "89",
	}
	err := rest.SaveAddress(address)
	if err != nil {
		err.Error()
	}
	score = database.Score{
		Category: &database.Category{
			Name: "Polka",
		},
		Difficulty: 3,
		Price:      39.9,
		Title:      "Eine letzte Runde (Blasorchesterfassung)",
	}
	if err != nil {
		err.Error()
	}
	err = rest.SaveScore(score)
	order = database.Order{
		BillingAddress:  address,
		Company:         "Millionen Show",
		Date:            1568024628000,
		DeliveryAddress: address,
		Email:           "jauch@werwirdswohl.de",
		FirstName:       "Günter",
		LastName:        "Jauch",
		Payed:           false,
		Salutation:      "Herr",
		Score:           score,
		ScoreAmount:     3,
		Telephone:       "+4954783907",
	}
	err = rest.SaveOrder(order)
	if err != nil {
		err.Error()
	}
}

//Inserts TestData to the Database
func insertTestData2() {
	address = database.Address{
		City:         "Wien",
		PostCode:     "1050",
		State:        "Österreich",
		Street:       "Spengergasse",
		StreetNumber: "20",
	}
	err := rest.SaveAddress(address)
	if err != nil {
		err.Error()
	}
	score = database.Score{
		Category: &database.Category{
			Name: "Polka",
		},
		Difficulty: 1,
		Price:      24.37,
		Title:      "Über den (Netzwerk)Brücken (Wiener Linien Fassung)",
	}
	err = rest.SaveScore(score)
	if err != nil {
		err.Error()
	}
	order = database.Order{
		BillingAddress:  address,
		Company:         "",
		Date:            10022018,
		DeliveryAddress: address,
		Email:           "hpberger@spengergasse.at",
		FirstName:       "Hans-Peter",
		LastName:        "Berger",
		Payed:           true,
		Salutation:      "Herr",
		Score:           score,
		ScoreAmount:     1,
		Telephone:       "",
	}
	err = rest.SaveOrder(order)
	if err != nil {
		err.Error()
	}
}

func insertTestData3() {
	address = database.Address{
		City:         "Leopoldsdorf im Marchfelde",
		PostCode:     "2285",
		State:        "Österreich",
		Street:       "Leopold-Figl-Gasse",
		StreetNumber: "2c",
	}
	err := rest.SaveAddress(address)
	if err != nil {
		err.Error()
	}
	score = database.Score{
		Category: &database.Category{
			Name: "Marsch",
		},
		Difficulty: 1,
		Price:      24.37,
		Title:      "Arnheim Marsch",
	}
	err = rest.SaveScore(score)
	if err != nil {
		err.Error()
	}
	order = database.Order{
		BillingAddress:  address,
		Company:         "",
		Date:            20190201,
		DeliveryAddress: address,
		Email:           "e11908080@student.tuwien.ac.at",
		FirstName:       "Richard",
		LastName:        "Stëckl",
		Payed:           true,
		Salutation:      "Herr",
		Score:           score,
		ScoreAmount:     1,
		Telephone:       "",
	}
	err = rest.SaveOrder(order)
	if err != nil {
		err.Error()
	}
}
