package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")

	img, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	kafkaResp, err := h.service.Process(img)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	resp, err := json.Marshal(kafkaResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)

	w.WriteHeader(http.StatusCreated)
}
