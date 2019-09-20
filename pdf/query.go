package pdf

import "github.com/mknentwich/core/database"

//Contains data from database for better handling
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
