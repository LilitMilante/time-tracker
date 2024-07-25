package tracker

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

type PassportNumber struct {
	PassportNumber string `json:"passportNumber"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req PassportNumber

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parts := strings.Split(req.PassportNumber, " ")

	if len(parts) != 2 {
		http.Error(w, "passport number mast has format '1234 567890'", http.StatusBadRequest)
		return
	}

	if len(parts[0]) != 4 {
		http.Error(w, "passport series must contains 4 signs", http.StatusBadRequest)
		return
	}

	passportSeries, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "passport series must be a integer", http.StatusBadRequest)
		return
	}

	if len(parts[1]) != 6 {
		http.Error(w, "passport number must contains 6 signs", http.StatusBadRequest)
		return
	}

	passportNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "passport number must be a integer", http.StatusBadRequest)
		return
	}

	err = h.s.CreateUser(ctx, passportSeries, passportNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
