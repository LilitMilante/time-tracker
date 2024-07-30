package tracker

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	_ "time-tracker/docs"

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

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with passport number
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			passportNumber	body		PassportNumber	true	"Passport number in format '1234 567890'"
//	@Success		200				{object}	User
//	@Failure		400				{string}	string	"Invalid input"
//	@Failure		500				{string}	string	"Internal error"
//	@Router			/users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

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
		l.Error("create user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateUser godoc
//
//	@Summary		Update an existing user
//	@Description	Update user details
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		UpdateUser	true	"User to update"
//	@Success		200		{object}	User
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		404		{string}	string	"User not found"
//	@Failure		500		{string}	string	"Internal error"
//	@Router			/users [patch]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	var updUser UpdateUser

	err := json.NewDecoder(r.Body).Decode(&updUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.s.UpdateUser(ctx, updUser)
	if err != nil {
		l.Error("update user", "error", err)
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

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Produce		json
//	@Param			user_id	path		string	true	"User ID"
//	@Success		200		{string}	string	"User deleted"
//	@Failure		400		{string}	string	"Invalid user ID"
//	@Failure		404		{string}	string	"User not found"
//	@Failure		500		{string}	string	"Internal error"
//	@Router			/users/{user_id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	id, err := uuid.FromString(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.DeleteUser(ctx, id)
	if err != nil {
		l.Error("delete user", "error", err)
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

// StartWork godoc
//
//	@Summary		Start work on a task
//	@Description	Start work on a task for a user
//	@Tags			work
//	@Accept			json
//	@Produce		json
//	@Param			startWorkRequest	body		StartWorkRequest	true	"Start work request"
//	@Success		200					{string}	string				"Work started"
//	@Failure		400					{string}	string				"Invalid input"
//	@Failure		500					{string}	string				"Internal error"
//	@Router			/work/start [post]
func (h *Handler) StartWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	var req StartWorkRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.StartWork(ctx, req.UserID, req.TaskID)
	if err != nil {
		l.Error("start work", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type FinishWorkRequest struct {
	UserID uuid.UUID `json:"user_id"`
	TaskID uuid.UUID `json:"task_id"`
}

// FinishWork godoc
//
//	@Summary		Finish work on a task
//	@Description	Finish work on a task for a user
//	@Tags			work
//	@Accept			json
//	@Produce		json
//	@Param			finishWorkRequest	body		FinishWorkRequest	true	"Finish work request"
//	@Success		200					{string}	string				"Work finished"
//	@Failure		400					{string}	string				"Invalid input"
//	@Failure		404					{string}	string				"Task not found"
//	@Failure		500					{string}	string				"Internal error"
//	@Router			/work/finish [post]
func (h *Handler) FinishWork(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	var req FinishWorkRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.FinishWork(ctx, req.UserID, req.TaskID)
	if err != nil {
		l.Error("finish work", "error", err)
		if errors.Is(err, ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TaskSpendTimesByUser godoc
//
//	@Summary		Get task spend times by user
//	@Description	Get the time spent on tasks by a user within a specified period
//	@Tags			tasks
//	@Produce		json
//	@Param			user_id		path		string	true	"User ID"
//	@Param			start_date	query		string	false	"Start date in format 'DD-MM-YYYY'"
//	@Param			end_date	query		string	false	"End date in format 'DD-MM-YYYY'"
//	@Success		200			{object}	[]TaskSpendTime
//	@Failure		400			{string}	string	"Invalid input"
//	@Failure		404			{string}	string	"User or task not found"
//	@Failure		500			{string}	string	"Internal error"
//	@Router			/users/{user_id}/report [get]
func (h *Handler) TaskSpendTimesByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)
	var err error

	id, err := uuid.FromString(r.PathValue("user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var period Period

	if startDate != "" {
		period.StartDate, err = time.Parse("02-01-2006", startDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if endDate != "" {
		period.EndDate, err = time.Parse("02-01-2006", endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if period.StartDate.Equal(period.EndDate) {
		// to get values by one day (with the same dates)
		period.EndDate = period.EndDate.Add(time.Hour * 24)
	}

	spendTimesByUser, err := h.s.TaskSpendTimesByUser(ctx, id, period)
	if err != nil {
		l.Error("get task spend times by user", "error", err)
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

// Users godoc
//
//	@Summary		Get users
//	@Description	Get a list of users with optional filters
//	@Tags			users
//	@Produce		json
//	@Param			page			query		int		false	"Page number"
//	@Param			per_page		query		int		false	"Number of users per page"
//	@Param			id				query		string	false	"User ID"
//	@Param			passport_series	query		int		false	"Passport series"
//	@Param			passport_number	query		int		false	"Passport number"
//	@Param			surname			query		string	false	"Surname"
//	@Param			name			query		string	false	"Name"
//	@Param			patronymic		query		string	false	"Patronymic"
//	@Param			address			query		string	false	"Address"
//	@Success		200				{object}	[]User
//	@Failure		400				{string}	string	"Invalid input"
//	@Failure		500				{string}	string	"Internal error"
//	@Router			/users [get]
func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)
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
		l.Error("get users", "errors", err)
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
