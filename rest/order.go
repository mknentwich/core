package rest

import (
	"github.com/mknentwich/core/database"
	"time"
)

type PostedOrder struct {
	AddressesEqual bool
	Bcity          string
	BpostCode      string
	Bstate         string
	Bstreet        string
	BstreetNumber  string
	City           string
	PostCode       string
	ScoreId        uint
	State          string
	Street         string
	StreetNumber   string
	Company        string
	Email          string
	FirstName      string
	LastName       string
	Salutation     string
	Telephone      string
}

func (p *PostedOrder) Order() *database.Order {
	address := &database.Address{
		City:         p.City,
		PostCode:     p.PostCode,
		State:        &database.State{
			Name: p.State,
		},
		Street:       p.Street,
		StreetNumber: p.StreetNumber,
	}
	var billingAddress *database.Address
	if p.AddressesEqual {
		billingAddress = address
	} else {
		billingAddress = &database.Address{
			City:         p.Bcity,
			PostCode:     p.BpostCode,
			State:        &database.State{
				Name: p.Bstate,
			},
			Street:       p.Bstreet,
			StreetNumber: p.BstreetNumber,
		}
	}
	now := time.Now()
	order := &database.Order{
		BillingAddress:  billingAddress,
		Company:         p.Company,
		Date:            now.Unix(),
		DeliveryAddress: address,
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
