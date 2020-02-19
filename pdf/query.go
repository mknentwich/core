package pdf

import (
	"errors"
	"github.com/mknentwich/core/database"
)

//Contains data for the PDF Entry
//the address is the deliveryAddress only
type OrderResultPDF struct {
	City           string
	PostCode       string
	Name           string
	DeliveryPrice  float64
	Street         string
	StreetNumber   string
	BillingAddress *BillingAddress
	Company        string
	Date           int
	FirstName      string
	LastName       string
	Salutation     string
	ScoreItems     []*ScoreItems
	ReferenceCount int
	BillingDate    int
}

//BillingAddress for the OrderResultPDF
type BillingAddress struct {
	City         string
	PostCode     string
	Street       string
	StreetNumber string
	Name         string
}

//Contains scores for the OrderResultPDF
type ScoreItems struct {
	Title       string
	Price       float64
	ScoreAmount int
}

//Selects order by ID and serves a result struct for better bill handling
//The billingAddress is queried separately, as the whole sql query can't be done in one part
func QueryOrderFromIdForPDF(id int) (OrderResultPDF, error) {
	var pdfOrderResult OrderResultPDF
	var billing BillingAddress
	var scoreItems []*ScoreItems
	recordNotFound := database.Receive().Table("orders").Select("city, post_code, street, street_number, name, delivery_price, "+
		"orders.id, orders.company, orders.date, orders.first_name, orders.last_name, orders.salutation, orders.reference_count, orders.billing_date").
		Joins("inner join addresses on orders.delivery_address_id = addresses.id").
		Joins("JOIN states on states.id = addresses.state_id").
		Where("orders.id = ?", id).
		Find(&pdfOrderResult).Scan(&pdfOrderResult).RecordNotFound()
	if recordNotFound {
		return pdfOrderResult, errors.New("QueryOrderFromIdForPDF: Record with orderID not found")
	}
	recordNotFound = database.Receive().Table("orders").Select("city, post_code, street, street_number, name").
		Joins("inner join addresses on orders.billing_address_id = addresses.id").
		Joins("JOIN states on states.id = addresses.state_id").
		Where("orders.id = ?", id).
		Find(&billing).Scan(&billing).RecordNotFound()
	if recordNotFound {
		return pdfOrderResult, errors.New("QueryOrderFromIdForPDF: Record with orderID not found")
	}
	recordNotFound = database.Receive().Table("items").Select("items.score_amount, scores.title, scores.price").
		Joins("inner join scores on items.score_id = scores.id").
		Where("items.order_id = ?", id).
		Find(&scoreItems).Scan(&scoreItems).RecordNotFound()
	if recordNotFound {
		return pdfOrderResult, errors.New("QueryOrderFromIdForPDF: Record with orderID not found")
	}
	pdfOrderResult.BillingAddress = &billing
	pdfOrderResult.ScoreItems = scoreItems
	return pdfOrderResult, nil
}
