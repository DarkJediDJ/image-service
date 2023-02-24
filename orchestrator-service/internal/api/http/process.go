package http_handler

import (
	"io"
	"net/http"
)

func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	img, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	proccessedImg, err := h.service.Process(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(proccessedImg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
