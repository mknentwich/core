package rest

import (
	"github.com/mknentwich/core/database"
	"time"
)

type PostedOrder struct {
	Dcity         string
	DpostCode     string
	Dstate        string
	Dstreet       string
	DstreetNumber string
	City          string
	PostCode      string
	ScoreId       uint
	State         string
	Street        string
	StreetNumber  string
	Company       string
	Email         string
	FirstName     string
	LastName      string
	Salutation    string
	Telephone     string
}

func (p *PostedOrder) Order() *database.Order {
	address := &database.Address{
		City:     p.City,
		PostCode: p.PostCode,
		State: &database.State{
			Name: p.State,
		},
		Street:       p.Street,
		StreetNumber: p.StreetNumber,
	}
	deliveryAddress := &database.Address{
		City:     p.Dcity,
		PostCode: p.DpostCode,
		State: &database.State{
			Name: p.Dstate,
		},
		Street:       p.Dstreet,
		StreetNumber: p.DstreetNumber,
	}
	now := time.Now()
	order := &database.Order{
		BillingAddress:  address,
		Company:         p.Company,
		Date:            now.Unix(),
		DeliveryAddress: deliveryAddress,
		Email:           p.Email,
		FirstName:       p.FirstName,
		LastName:        p.LastName,
		Payed:           false,
		ReferenceCount:  0,
		Salutation:      p.Salutation,
		ScoreID:         p.ScoreId,
		ScoreAmount:     1,
		Telephone:       p.Telephone,
	}
	return order
}
