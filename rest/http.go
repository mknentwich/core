package rest

import (
	"encoding/json"
	"github.com/mknentwich/core/context"
	"github.com/mknentwich/core/utils"
	"net/http"
	"strconv"
	"strings"
)

//Logging function for this package.
var log context.Log

//Serve function for this package.
func Serve(args context.ServiceArguments) (context.ServiceResult, error) {
	log = args.Log
	mux := http.NewServeMux()
	mux.HandleFunc("/", utils.HttpImplement(log))
	mux.HandleFunc("/categories", utils.Rest(flat(get(QueryCategoriesFlat), get(QueryCategoriesWithChildrenAndScores))))
	mux.HandleFunc("/order", utils.Cors(postOrder))
	mux.HandleFunc("/scores", utils.Rest(get(QueryScoresFlat)))
	return context.ServiceResult{HttpHandler: mux}, initializeTemplates()
}

//Encodes a structure as JSON and returns it
func get(query DataQuery) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			encoder := json.NewEncoder(writer)
			err := encoder.Encode(query())
			if err != nil {
				log(context.LOG_ERROR, "An error occurred on return a REST GET request: %s", err.Error())
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

//calls flat handler, if url attribute `flat` is `true`, otherwise is call treeHandler.
func flat(flatHandler http.HandlerFunc, treeHandler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("flat") == "true" {
			flatHandler(rw, r)
		} else {
			treeHandler(rw, r)
		}
	}
}

func postOrder(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	f := r.PostForm
	scoreId, err := strconv.ParseUint(f.Get("scoreId"), 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	o := PostedOrder{
		AddressesEqual: strings.ToLower(f.Get("addressesEqual")) == "on",
		Bcity:          f.Get("bcity"),
		BpostCode:      f.Get("bpostcode"),
		Bstate:         f.Get("bstate"),
		Bstreet:        f.Get("bstreet"),
		BstreetNumber:  f.Get("bstreetNumber"),
		City:           f.Get("city"),
		PostCode:       f.Get("postCode"),
		ScoreId:        uint(scoreId),
		State:          f.Get("state"),
		Street:         f.Get("street"),
		StreetNumber:   f.Get("streetNumber"),
		Company:        f.Get("company"),
		Email:          f.Get("email"),
		FirstName:      f.Get("firstName"),
		LastName:       f.Get("lastName"),
		Salutation:     f.Get("salutation"),
		Telephone:      f.Get("telephone"),
	}
	order := (&o).Order()
	err = InsertNewOrder(order)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log(context.LOG_ERROR, "error occurred while persisting a new order: %s", err.Error())
	}
	err = notify(order)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log(context.LOG_ERROR, "error occurred while sending order mails: %s", err.Error())
	}
}
