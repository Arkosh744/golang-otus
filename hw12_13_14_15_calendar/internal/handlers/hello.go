package handlers

import (
	"io"
	"net/http"
)

func (h *Handler) hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello, HTTP!\n")
}
