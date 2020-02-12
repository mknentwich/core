package rest

import (
	"github.com/jinzhu/gorm"
	"github.com/mknentwich/core/database"
)

type DataQuery func() interface{}

func QueryCategoriesWithChildrenAndScores() interface{} {
	return QueryCategoriesWithChildrenAndScoresPreserve()
}

func QueryCategoriesWithChildrenAndScoresPreserve() []database.Category {
	categories := make([]database.Category, 0)
	database.Receive().Preload("Children", func(d *gorm.DB) *gorm.DB {
		return d.Order("name").Preload("Scores", func(d2 *gorm.DB) *gorm.DB {
			return d2.Order("title")
		})
	}).Where("categories.parent_id is null").Order("name").Find(&categories)
	return categories
}

func QueryOrders(payed bool) []database.Order {
	orders := make([]database.Order, 0)
	database.Receive().Where("payed = ?", payed).Find(&orders)
	return orders
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

//Inserts or updates a address to the database
func SaveAddress(address database.Address) error {
	err := database.Receive().Save(&address).Error
	return err
}

//Inserts or updates a category to the database
func SaveCategory(category database.Category) error {
	err := database.Receive().Save(&category).Error
	return err
}

//Inserts or updates a score to the database
func SaveScore(score database.Score) error {
	err := database.Receive().Save(&score).Error
	return err
}

//Inserts new order to the database
func InsertNewOrder(order *database.Order) error {
	err := database.Receive().Create(order).Error
	database.Receive().Preload("Score").Find(order)
	return err
}

//Inserts or updates a order to the database
func SaveOrder(order database.Order) error {
	err := database.Receive().Save(&order).Error
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
