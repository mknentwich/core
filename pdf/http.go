package pdf

import (
	"fmt"
	"github.com/mknentwich/core/auth"
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
	//TODO wrap order with `Auth`
	mux.HandleFunc("/order/", auth.Auth(pdfHeader(httpOrder)))
	return context.ServiceResult{HttpHandler: mux}, nil
}

//Http handler for bills
func httpOrder(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	idx := strings.LastIndex(path, "/") + 1
	orderId, err := strconv.Atoi(path[idx:])
	if err != nil {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	log(context.LOG_INFO, "requested order with id: %d", orderId)
	//TODO define filename and append some infos such as the number of the bill
	filename := "des-is-mei-rechnung.pdf"
	rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	//TODO write pdf stream to response body
	utils.HttpImplement(log)(rw, r)
}

//Sets the correct PDF HTTP headers
func pdfHeader(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/pdf")
		handlerFunc(rw, r)
	}
}
