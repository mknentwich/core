package rest

import (
	"github.com/jinzhu/gorm"
	"github.com/mknentwich/core/database"
)

type DataQuery func() interface{}

type OrderResultPDF struct {
	City           string
	PostCode       string
	State          string
	Street         string
	StreetNumber   string
	ID             uint
	Company        string
	Date           int
	FirstName      string
	LastName       string
	Salutation     string
	ScoreAmount    int
	Title          string
	Price          float64
	ReferenceCount int
}

func QueryCategoriesWithChildrenAndScores() interface{} {
	categories := make([]database.Category, 0)
	database.Receive().Preload("Children", func(d *gorm.DB) *gorm.DB {
		return d.Order("name").Preload("Scores", func(d2 *gorm.DB) *gorm.DB {
			return d2.Order("title")
		})
	}).Where("categories.parent_id is null").Order("name").Find(&categories)
	return categories
}

//Selects order by ID and serves a result struct for better bill handling
func QueryOrderFromIdForPDF(id int) OrderResultPDF {
	var pdfOrderResult OrderResultPDF
	database.Receive().Table("orders").Select("addresses.city, addresses.post_code, addresses.state, addresses.street, addresses.street_number, "+
		"orders.id, orders.company, orders.date, orders.first_name, orders.last_name, orders.salutation, orders.score_amount, orders.reference_count, scores.title, scores.price").Joins("inner join addresses on orders.billing_address_id = addresses.id").
		Joins("inner join scores on orders.score_id = scores.id").Where("orders.id = ?", id).Find(&pdfOrderResult).Scan(&pdfOrderResult)
	return pdfOrderResult
}

//Finds the highest referenceCount value from the orders table
func FindMaxReferenceCount() int {
	var max int
	row := database.Receive().Table("orders").Select("MAX(reference_count)").Row()
	row.Scan(&max)
	return max
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

//Queries all scores without a tree structure.
func QueryScoresFlat() interface{} {
	scores := make([]database.Score, 0)
	database.Receive().Find(&scores)
	mp := make(map[uint]database.Score)
	for _, score := range scores {
		mp[score.ID] = score
	}
	return mp
}

//Queries all categories without a tree structure.
func QueryCategoriesFlat() interface{} {
	categories := make([]database.Category, 0)
	database.Receive().Find(&categories)
	mp := make(map[uint]database.Category)
	for _, category := range categories {
		mp[category.ID] = category
	}
	return mp
}
