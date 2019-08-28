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
