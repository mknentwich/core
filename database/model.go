package database

import "github.com/jinzhu/gorm"

type Address struct {
	gorm.Model
	City         string
	PostCode     string
	State        string
	Street       string
	StreetNumber string
}

type Category struct {
	gorm.Model
	Children []Category `gorm:"foreignkey:ParentID"`
	Name     string
	Parent   *Category `gorm:"PRELOAD:false"`
	ParentID uint      `sql:"type:integer REFERENCES categories(id)"`
	Scores   []Score   `gorm:"foreignkey:CategoryID"`
}

type Score struct {
	gorm.Model
	Category   *Category
	CategoryID uint `sql:"type:integer REFERENCES categories(id)"`
	Difficulty int
	Price      float64
	Title      string
}

type Order struct {
	gorm.Model
	BillingAddress    Address
	BillingAddressID  uint `sql:"type:integer REFERENCES addresses(id)"`
	Company           string
	Date              int
	DeliveryAddress   Address
	DeliveryAddressID uint `sql:"type:integer REFERENCES addresses(id)"`
	Email             string
	FirstName         string
	LastName          string
	Payed             bool
	ReferenceCount    int
	Salutation        string
	Score             Score
	ScoreID           uint `sql:"type:integer REFERENCES scores(id)"`
	ScoreAmount       int
	Telephone         string
}

type User struct {
	gorm.Model
	Email      string `gorm:"primary_key"`
	Admin      bool
	Name       string
	Password   string
	LastChange int
	LastLogin  int
}
