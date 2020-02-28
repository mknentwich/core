package database

import (
	"testing"
	"time"
)

func TestInitializeDb(t *testing.T) {
	err := initializeDb()
	if err != nil {
		t.Errorf("Error on creating the database: %s", err.Error())
	}
}

func TestAddStateTrigger(t *testing.T) {
	err := initializeDb()
	if err != nil {
		t.Errorf("Error on creating the database: %s", err.Error())
	}
	var state = State{
		Name: "Österreich",
	}
	var billing = Address{
		City:         "Wien",
		PostCode:     "1010",
		State:        &state,
		Street:       "Testbillingstraße",
		StreetNumber: "10",
	}
	var delivery = Address{
		City:         "Wien",
		PostCode:     "1010",
		State:        &state,
		Street:       "Testdeliverystraße",
		StreetNumber: "20",
	}
	var order = Order{
		BillingAddress:  &billing,
		Company:         "MVL",
		Date:            time.Now().Unix(),
		DeliveryAddress: &delivery,
		Email:           "test@test.at",
		FirstName:       "Max",
		LastName:        "Musater",
		Payed:           false,
		Salutation:      "Herr",
		Score: Score{
			Category: &Category{
				Name: "Polka",
			},
			Difficulty: 2,
			Price:      39.99,
			Title:      "Eine letzte Runde",
		},
		ScoreAmount: 1,
		Telephone:   "+136565406",
	}
	err = Receive().Create(&order).Error

	country := State{
		Name: "Österreich",
	}
	address := Address{
		City:         "Leopoldsdorf im Marchfelde",
		PostCode:     "2285",
		State:        &country,
		Street:       "Leopold-Figl-Gasse",
		StreetNumber: "2c",
	}
	score := Score{
		Category: &Category{
			Name: "Marsch",
		},
		Difficulty: 1,
		Price:      24.37,
		Title:      "Arnheim Marsch",
	}
	order2 := Order{
		BillingAddress: &Address{
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
		Score:           score,
		ScoreAmount:     2,
		Telephone:       "",
	}
	err = Receive().Create(&order2).Error
	var count = 0
	Receive().Table("states").Where("name = ?", "Österreich").Count(&count)
	if count != 1 {
		t.Errorf("%s", "State was inserted into database more than one time!")
	}
}
