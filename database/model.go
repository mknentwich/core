package database

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Address struct {
	gorm.Model
	City         string `json:"city"`
	PostCode     string `json:"postCode"`
	State        *State
	StateID      *uint  `sql:"type:integer REFERENCES states(id)" json:"-"`
	Street       string `json:"street"`
	StreetNumber string `json:"streetNumber"`
}

type State struct {
	gorm.Model
	Name          string  `json:"name"`
	DeliveryPrice float64 `json:"deliveryPrice"`
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
	BillingAddress    *Address `json:"billingAddress"`
	BillingAddressID  uint     `sql:"type:integer REFERENCES addresses(id)" json:"-"`
	BillingDate       int64    `json:"billingDate"`
	Company           string   `json:"company"`
	Date              int64    `json:"date"`
	DeliveryAddress   *Address `json:"deliveryAddress"`
	DeliveryAddressID uint     `sql:"type:integer REFERENCES addresses(id)" json:"-"`
	Email             string   `json:"email"`
	FirstName         string   `json:"firstName"`
	LastName          string   `json:"lastName"`
	Payed             bool     `json:"payed"`
	ReferenceCount    int      `json:"referenceCount"`
	Salutation        string   `json:"salutation"`
	Score             Score    `json:"score"`
	ScoreID           uint     `sql:"type:integer REFERENCES scores(id)" json:"scoreId"`
	ScoreAmount       int      `json:"scoreAmount"`
	Telephone         string   `json:"telephone"`
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

//Hook for generating BillingDate and ReferenceCount before Order is saved to the database
func (order *Order) BeforeSave(db *gorm.DB) (err error) {
	if order.ReferenceCount == 0 && order.BillingDate == 0 {
		maxRef := 0
		now := time.Now()
		fTime, err := strconv.Atoi(now.Format("20060102"))
		if err != nil {
			return err
		}
		var s sql.NullString
		row := db.Table("orders").Select("max(reference_count)").Where("billing_date = ?", &fTime).Row()
		err = row.Scan(&s)
		if err != nil {
			return err
		}
		if s.Valid {
			maxRef, err = strconv.Atoi(s.String)
			if err != nil {
				return err
			}
		}
		maxRef++
		order.BillingDate = int64(fTime)
		order.ReferenceCount = maxRef
	}
	return
}
