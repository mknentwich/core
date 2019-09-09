package rest

import (
	"github.com/jinzhu/gorm"
	"github.com/mknentwich/core/database"
)

type DataQuery func() interface{}

func QueryCategoriesWithChildrenAndScores() interface{} {
	categories := make([]database.Category, 0)
	database.Receive().Preload("Children", func(d *gorm.DB) *gorm.DB {
		return d.Joins("inner join scores s on categories.id = s.category_id").Order("name").Preload("Scores", func(d2 *gorm.DB) *gorm.DB {
			return d2.Order("title")
		})
	}).Where("categories.parent_id is null").Joins("inner join categories c on c.parent_id = categories.id").Joins("inner join scores s on s.category_id = c.id").Order("name").Find(&categories)
	return categories
}

type PDFOrderResult struct {
	City         string
	PostCode     string
	State        string
	Street       string
	StreetNumber string
	ID           uint
	Company      string
	Date         int
	FirstName    string
	LastName     string
	Salutation   string
	ScoreAmount  int
	Title        string
	Price        float64
}

func QueryOrderFromIdForPDF(id int) PDFOrderResult {
	var pdfOrderResult PDFOrderResult
	database.Receive().Table("orders").Select("addresses.city, addresses.post_code, addresses.state, addresses.street, addresses.street_number, "+
		"orders.id, orders.company, orders.date, orders.first_name, orders.last_name, orders.salutation, orders.score_amount, scores.title, scores.price").Joins("inner join addresses on orders.billing_address_id = addresses.id").
		Joins("inner join scores on orders.score_id = scores.id").Where("orders.id = ?", id).Find(&pdfOrderResult).Scan(&pdfOrderResult)
	return pdfOrderResult
}

//Inserts new address to the database
func InsertNewAddress(address database.Address) error {
	err := database.Receive().Create(&address).Error
	return err
}

//Inserts new category to the database
func InsertNewCategory(category database.Category) error {
	err := database.Receive().Create(&category).Error
	return err
}

//Inserts new score to the database
func InsertNewScore(score database.Score) error {
	err := database.Receive().Create(&score).Error
	return err
}

//Inserts new order to the database
func InsertNewOrder(order database.Order) error {
	err := database.Receive().Create(&order).Error
	return err
}
