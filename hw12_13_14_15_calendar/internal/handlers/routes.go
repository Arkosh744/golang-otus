package handlers

import (
	"net/http"
)

type CalendarService interface {
}

func InitRouter(serv CalendarService) *http.ServeMux {
	mux := http.NewServeMux()

	//addToCart := add.NewHandler(serv).Handle
	//mux.Handle("/addToCart", wrappers.New(addToCart))

	return mux
}
