package handlers

import (
	"io"
	"net/http"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/log"
)

func (h *Handler) hello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	if _, err := io.WriteString(w, "Hello, HTTP!\n"); err != nil {
		log.Errorf("failed to write response: %v", err)

		return
	}
}
