package dav

import "github.com/mknentwich/core/database"

func queryFilteredBills(filter billFilter) []database.Order {
	return make([]database.Order, 0)
}
