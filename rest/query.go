package rest

import (
	"github.com/jinzhu/gorm"
	"github.com/mknentwich/core/database"
)

type DataQuery func() interface{}

func QueryCategoriesWithChildrenAndScores() interface{} {
	categories := make([]database.Category, 0)
	database.Receive().Preload("Children", func(d *gorm.DB) *gorm.DB {
		return d.Order("name").Preload("Scores", func(d2 *gorm.DB) *gorm.DB {
			return d2.Order("title")
		})
	}).Where("parent_id is null").Order("name").Find(&categories)
	return categories
}
