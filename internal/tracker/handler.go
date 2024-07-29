package tracker

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
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

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updUser UpdateUser

	err := json.NewDecoder(r.Body).Decode(&updUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.s.UpdateUser(ctx, updUser)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.FromString(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type StartWorkRequest struct {
	UserID uuid.UUID `json:"user_id"`
	TaskID uuid.UUID `json:"task_id"`
}

func (h *Handler) StartWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req StartWorkRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.StartWork(ctx, req.UserID, req.TaskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type FinishWorkRequest struct {
	UserID uuid.UUID `json:"user_id"`
	TaskID uuid.UUID `json:"task_id"`
}

func (h *Handler) FinishWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req FinishWorkRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.FinishWork(ctx, req.UserID, req.TaskID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) TaskSpendTimesByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.FromString(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	spendTimesByUser, err := h.s.TaskSpendTimesByUser(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(spendTimesByUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	pageParam := r.URL.Query().Get("page")
	perPageParam := r.URL.Query().Get("per_page")

	page := 1
	perPage := 20

	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if perPageParam != "" {
		perPage, err = strconv.Atoi(perPageParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	filter, err := parsUserFilter(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := h.s.Users(ctx, page, perPage, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parsUserFilter(v url.Values) (f UserFilter, err error) {
	idParam := v.Get("id")
	if idParam != "" {
		id, err := uuid.FromString(idParam)
		if err != nil {
			return UserFilter{}, err
		}
		f.ID = &id
	}

	pSeries := v.Get("passport_series")
	if pSeries != "" {
		series, err := strconv.Atoi(pSeries)
		if err != nil {
			return UserFilter{}, err
		}
		f.PassportSeries = &series
	}

	pNumber := v.Get("passport_number")
	if pNumber != "" {
		number, err := strconv.Atoi(pNumber)
		if err != nil {
			return UserFilter{}, err
		}
		f.PassportNumber = &number
	}

	surname := v.Get("surname")
	if surname != "" {
		f.Surname = &surname
	}

	name := v.Get("name")
	if name != "" {
		f.Name = &name
	}

	patronymic := v.Get("patronymic")
	if patronymic != "" {
		f.Patronymic = &patronymic
	}

	address := v.Get("address")
	if address != "" {
		f.Address = &address
	}

	return f, nil
}
