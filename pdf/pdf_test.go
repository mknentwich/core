package pdf

import (
	"errors"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/rest"
	"io"
	"net/http"
	"os"
	"testing"
)

var country database.State
var address database.Address
var order database.Order
var score1 database.Score
var score2 database.Score
var scoreItem1 database.Item
var scoreItem2 database.Item

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
		reader, filename, err := GeneratePDF(1)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
			return
		}
		f, err := os.OpenFile("Rechnung_"+filename+".pdf", os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
		io.Copy(f, reader)
		reader, filename, err = GeneratePDF(2)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
			return
		}
		f, err = os.OpenFile("Rechnung_"+filename+".pdf", os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
		io.Copy(f, reader)
		reader, filename, err = GeneratePDF(3)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
			return
		}
		f, err = os.OpenFile("Rechnung_"+filename+".pdf", os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			t.Errorf("Error on creating the bill pdf: %s", err.Error())
		}
		io.Copy(f, reader)
	})
}

//Inserts TestData to the Database
func insertTestData() {
	country = database.State{
		Name:          "Deutschland",
		DeliveryPrice: 7,
	}
	address = database.Address{
		City:         "Leopoldsdorf im wundersc Marchfelde",
		PostCode:     "50354",
		State:        &country,
		Street:       "vTg>726X$Do5:x,Yt?qvBh#~Fl'Fy9",
		StreetNumber: "xx-5ax1",
	}
	err := rest.SaveAddress(address)
	if err != nil {
		err.Error()
	}
	score1 = database.Score{
		Category: &database.Category{
			Name: "Polka",
		},
		Difficulty: 3,
		Price:      39.9,
		Title:      "Eine letzte Runde (Blasorchesterfassung)",
	}
	score2 = database.Score{
		Category: &database.Category{
			Name: "Märsche",
		},
		Difficulty: 4,
		Price:      19.9,
		Title:      "Florentiner Marsch",
	}
	err = rest.SaveScore(score1)
	if err != nil {
		err.Error()
	}
	err = rest.SaveScore(score2)
	if err != nil {
		err.Error()
	}
	scoreItem1 = database.Item{
		Score:       &score1,
		ScoreAmount: 2,
	}
	scoreItem2 = database.Item{
		Score:       &score2,
		ScoreAmount: 1,
	}
	order = database.Order{
		BillingAddress:  &address,
		Company:         "vTg>726X$Do5:x,Yt?qvBh#~Fl'Fy9bd^SJ",
		Date:            1568024628000,
		DeliveryAddress: &address,
		Email:           "jauch@werwirdswohl.de",
		FirstName:       "Günter",
		LastName:        "Jauch",
		Payed:           false,
		Salutation:      "Herr",
		Items:           []*database.Item{&scoreItem1, &scoreItem2},
		Telephone:       "+4954783907",
	}
	err = rest.SaveOrder(order)
	if err != nil {
		err.Error()
	}
}

//Inserts TestData to the Database
func insertTestData2() {
	country = database.State{
		Name:          "Österreich",
		DeliveryPrice: 3,
	}
	address = database.Address{
		City:         "Wien",
		PostCode:     "1050",
		State:        &country,
		Street:       "Spengergasse",
		StreetNumber: "20",
	}
	err := rest.SaveAddress(address)
	if err != nil {
		err.Error()
	}
	score1 = database.Score{
		Category: &database.Category{
			Name: "Polka",
		},
		Difficulty: 1,
		Price:      24.37,
		Title:      "Über den (Netzwerk)Brücken (Wiener Linien Fassung)",
	}
	score2 = database.Score{
		Category: &database.Category{
			Name: "Märsche",
		},
		Difficulty: 4,
		Price:      19.9,
		Title:      "Florentiner Marsch",
	}
	err = rest.SaveScore(score1)
	if err != nil {
		err.Error()
	}
	err = rest.SaveScore(score2)
	if err != nil {
		err.Error()
	}
	scoreItem1 = database.Item{
		Score:       &score1,
		ScoreAmount: 1,
	}
	scoreItem2 = database.Item{
		Score:       &score2,
		ScoreAmount: 1,
	}
	order = database.Order{
		BillingAddress: &database.Address{
			City:     "",
			PostCode: "",
			State: &database.State{
				Name:          "",
				DeliveryPrice: 0,
			},
			Street:       "",
			StreetNumber: "",
		},
		Company:         "",
		Date:            10022018,
		DeliveryAddress: &address,
		Email:           "hpberger@spengergasse.at",
		FirstName:       "Hans-Peter",
		LastName:        "Berger",
		Payed:           true,
		Salutation:      "Herr",
		Items:           []*database.Item{&scoreItem1, &scoreItem2},
		Telephone:       "",
	}
	err = rest.SaveOrder(order)
	if err != nil {
		err.Error()
	}
}

func insertTestData3() {
	country = database.State{
		Name:          "Österreich",
		DeliveryPrice: 3,
	}
	address = database.Address{
		City:         "Leopoldsdorf im Marchfelde",
		PostCode:     "2285",
		State:        &country,
		Street:       "Leopold-Figl-Gasse",
		StreetNumber: "2c",
	}
	err := rest.SaveAddress(address)
	if err != nil {
		err.Error()
	}
	score1 = database.Score{
		Category: &database.Category{
			Name: "Marsch",
		},
		Difficulty: 1,
		Price:      24.37,
		Title:      "Arnheim Marsch",
	}
	score2 = database.Score{
		Category: &database.Category{
			Name: "Märsche",
		},
		Difficulty: 4,
		Price:      19.9,
		Title:      "Florentiner Marsch",
	}
	err = rest.SaveScore(score1)
	if err != nil {
		err.Error()
	}
	err = rest.SaveScore(score2)
	if err != nil {
		err.Error()
	}
	scoreItem1 = database.Item{
		Score:       &score1,
		ScoreAmount: 3,
	}
	scoreItem2 = database.Item{
		Score:       &score2,
		ScoreAmount: 4,
	}
	order = database.Order{
		BillingAddress: &database.Address{
			City:         "Leopoldsdorf im Marchfelde",
			PostCode:     "2285",
			State:        &country,
			Street:       "Kempfendorf",
			StreetNumber: "2",
		},
		Company:         "Musikverein Leopoldsdorf",
		Date:            20190201,
		DeliveryAddress: &address,
		Email:           "e11908080@student.tuwien.ac.at",
		FirstName:       "Richard",
		LastName:        "Stöckl",
		Payed:           true,
		Salutation:      "Herr",
		Items:           []*database.Item{&scoreItem1, &scoreItem2},
		Telephone:       "",
	}
	err = rest.SaveOrder(order)
	if err != nil {
		err.Error()
	}
}
