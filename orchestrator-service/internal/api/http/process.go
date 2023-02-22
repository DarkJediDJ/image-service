package http_handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")

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

	resp, err := json.Marshal(proccessedImg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}
