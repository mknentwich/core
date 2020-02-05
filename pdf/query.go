package pdf

import (
	"database/sql"
	"errors"
	"github.com/mknentwich/core/database"
	"strconv"
)

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
	BillingDate    int
}

//Selects order by ID and serves a result struct for better bill handling
func QueryOrderFromIdForPDF(id int) (OrderResultPDF, error) {
	var pdfOrderResult OrderResultPDF
	recordNotFound := database.Receive().Table("orders").Select("addresses.city, addresses.post_code, addresses.state, addresses.street, addresses.street_number, "+
		"orders.id, orders.company, orders.date, orders.first_name, orders.last_name, orders.salutation, orders.score_amount, orders.reference_count, orders.billing_date, scores.title, scores.price").Joins("inner join addresses on orders.billing_address_id = addresses.id").
		Joins("inner join scores on orders.score_id = scores.id").Where("orders.id = ?", id).Find(&pdfOrderResult).Scan(&pdfOrderResult).RecordNotFound()
	if recordNotFound {
		return pdfOrderResult, errors.New("QueryOrderFromIdForPDF: Record with orderID not found")
	}
	return pdfOrderResult, nil
}

//Selects the highest referenceCount of today
//Returns 0 if no referenceCount was set this day
func QueryMaxReferenceCountToday(date int) (int, error) {
	var s sql.NullString
	row := database.Receive().Table("orders").Select("max(reference_count)").Where("billing_date = ?", &date).Row()
	err := row.Scan(&s)
	if err != nil {
		return 0, err
	}
	if s.Valid {
		maxRef, err := strconv.Atoi(s.String)
		if err != nil {
			return 0, err
		}
		return maxRef, nil
	} else {
		return 0, nil
	}
}

//Updates given fields to orders. Used to save OrderResultStruct
func UpdateOrders(id int, order *database.Order) error {
	err := database.Receive().Table("orders").Where("orders.id = ?", id).Update(&order).Error
	return err
}
