package handlers

type Handler struct {
	service CalendarService
}

func NewHandlers(service CalendarService) *Handler {
	return &Handler{service: service}
}
