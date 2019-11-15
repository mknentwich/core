package database

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Address struct {
	gorm.Model
	City         string `json:"city"`
	PostCode     string `json:"postCode"`
	State        string `json:"state"`
	Street       string `json:"street"`
	StreetNumber string `json:"streetNumber"`
}

type Category struct {
	gorm.Model
	Children []Category `gorm:"foreignkey:ParentID" json:"children"`
	Name     string     `json:"name"`
	Parent   *Category
	ParentID *uint   `sql:"type:integer REFERENCES categories(id)" json:"-"`
	Scores   []Score `gorm:"foreignkey:CategoryID" json:"scores"`
}

type Score struct {
	gorm.Model
	Category   *Category
	CategoryID uint    `sql:"type:integer REFERENCES categories(id)" json:"-"`
	Difficulty int     `json:"difficulty"`
	Price      float64 `json:"price"`
	Title      string  `json:"title"`
}

type Order struct {
	gorm.Model
	BillingAddress    *Address   `json:"billingAddress"`
	BillingAddressID  uint       `sql:"type:integer REFERENCES addresses(id)" json:"-"`
	Company           string     `json:"company"`
	Date              *time.Time `json:"date"`
	DeliveryAddress   *Address   `json:"deliveryAddress"`
	DeliveryAddressID uint       `sql:"type:integer REFERENCES addresses(id)" json:"-"`
	Email             string     `json:"email"`
	FirstName         string     `json:"firstName"`
	LastName          string     `json:"lastName"`
	Payed             bool       `json:"payed"`
	ReferenceCount    int        `json:"referenceCount"`
	Salutation        string     `json:"salutation"`
	Score             Score      `json:"score"`
	ScoreID           uint       `sql:"type:integer REFERENCES scores(id)" json:"scoreId"`
	ScoreAmount       int        `json:"scoreAmount"`
	Telephone         string     `json:"telephone"`
}

type User struct {
	*UserWithoutPassword
	Password string `json:"password"`
}

type UserWithoutPassword struct {
	gorm.Model
	Email      string `gorm:"primary_key" json:"email"`
	Admin      bool   `json:"admin"`
	Name       string `json:"name"`
	LastChange int    `json:"lastChange"`
	LastLogin  int    `json:"lastLogin"`
}
