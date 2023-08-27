package handlers

import (
	"net/http"
)

type CalendarService interface{}

func InitRouter(serv CalendarService) *http.ServeMux {
	mux := http.NewServeMux()
	handlers := NewHandlers(serv)

	mux.HandleFunc("/hello", handlers.hello)

	return mux
}
