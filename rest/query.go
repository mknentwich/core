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

//Inserts or updates a order to the database
func SaveOrder(order database.Order) error {
	err := database.Receive().Save(&order).Error
	return err
}
