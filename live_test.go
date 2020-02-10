package core

import (
	"errors"
	"github.com/mknentwich/core/auth"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/database"
	"github.com/mknentwich/core/dav"
	"github.com/mknentwich/core/media"
	"github.com/mknentwich/core/pdf"
	"github.com/mknentwich/core/rest"
	"github.com/mknentwich/core/template"
	"net/http"
	"net/mail"
	"testing"
	"time"
)

func InsertTestData() {
	andulko := database.Score{
		Title: "Andulko Safarova",
		Price: 40}
	andulko.ID = 73
	categories := []database.Category{
		{
			Name: "Märsche",
			Children: []database.Category{
				{
					Name: "Behmische Märsche",
					Scores: []database.Score{
						{
							Title: "Koline Koline",
							Price: 50},
						andulko,
						{
							Title: "Castaldo",
							Price: 74}},
				},
				{
					Name: "Österreichische Märsche",
					Scores: []database.Score{
						{
							Title: "Hoch und Deutschmeister",
							Price: 98},
						{
							Title: "Alte Kameraden",
							Price: 36}}}},
			Scores: []database.Score{
				{
					Title: "Arnheim",
					Price: 88}}},
		{
			Name: "Polka",
			Scores: []database.Score{
				{
					Title: "Ma Posledni",
					Price: 40},
				{
					Title: "Gute Nacht",
					Price: 33}},
		}}
	adr := database.Address{
		City:         "Hürth",
		PostCode:     "50354",
		State:        "Deutschland",
		Street:       "Kalscheurener Straße",
		StreetNumber: "89",
	}
	t, _ := time.Parse(time.RFC822, "01 Jan 15 10:00 UTC")
	orders := []database.Order{
		{
			BillingAddress:  &adr,
			Company:         "Millionen Show",
			Date:            &t,
			DeliveryAddress: &adr,
			Email:           "jauch@werwirdswohl.de",
			FirstName:       "Günter",
			LastName:        "Jauch",
			Payed:           false,
			Salutation:      "Herr",
			Score: database.Score{
				Difficulty: 3,
				Price:      39.9,
				Title:      "Eine letzte Runde (Blasorchesterfassung)",
				Category: &database.Category{
					Name: "Polka",
				},
			},
			ScoreAmount: 3,
			Telephone:   "",
		},
	}
	for _, v := range orders {
		rest.SaveOrder(v)
	}
	for _, category := range categories {
		rest.SaveCategory(category)
	}
	auth.SaveUser(&database.User{
		UserWithoutPassword: &database.UserWithoutPassword{
			Email: "albert@gmx.at",
			Admin: true,
			Name:  "Albert"},
		Password: "123456",
	})
}

func TestLive(t *testing.T) {
	conf := context.Configuration{
		GeneratedDirectory:   "gen",
		Authentication:       false,
		Host:                 "127.0.0.1:9400",
		JWTExpirationMinutes: 4,
		JWTSecret:            "9ef9486cf0a0e0ed17c2daa34a1e35f7",
		SQLiteFile:           "baum.db",
		Mail: context.EmailCredentials{
			Username: "",
			Password: "",
			SMTPHost: "",
			Address: &mail.Address{
				Name:    "",
				Address: "",
			},
		},
		OrderRetrievers: []*mail.Address{
			{
				Name:    "",
				Address: "",
			},
			{
				Name:    "",
				Address: ""}}}
	services := map[string]context.Serve{
		"db":       database.Serve,
		"dav":      dav.Serve,
		"api":      rest.Serve,
		"auth":     auth.Serve,
		"media":    media.Serve,
		"template": template.Serve,
		"pdf":      pdf.Serve}
	err := make(chan error)
	go func() {
		err <- context.InitializeCustomConfig(services, &conf)
	}()
	for err := errors.New(""); err != nil; _, err = http.Get("http://127.0.0.1:9400/") {
	}
	database.Receive().LogMode(true)
	InsertTestData()
	if e := <-err; e != nil {
		t.Error(e.Error())
	}
}
