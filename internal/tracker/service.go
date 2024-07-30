package tracker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

var ErrNotFound = errors.New("not found")
var ErrWorkAlreadyStarted = errors.New("work already started")

type Service struct {
	repo   *Repository
	client *http.Client
	apiURL string
}

func NewService(repo *Repository, apiURL string) *Service {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	return &Service{
		repo:   repo,
		client: client,
		apiURL: apiURL,
	}
}

func (s *Service) CreateUser(ctx context.Context, passportSeries int, passportNumber int) error {
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	l.Debug("get user info...")
	user, err := s.getUserInfo(ctx, passportSeries, passportNumber)
	if err != nil {
		return fmt.Errorf("get user info: %w", err)
	}

	user.ID = uuid.Must(uuid.NewV4())
	user.PassportSeries = passportSeries
	user.PassportNumber = passportNumber
	user.CreatedAt = time.Now()

	l.Debug("create user...")
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Service) getUserInfo(ctx context.Context, passportSeries, passportNumber int) (User, error) {
	url := fmt.Sprintf("%s/info?passportSerie=%d&passportNumber=%d", s.apiURL, passportSeries, passportNumber)

	var req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return User{}, fmt.Errorf("create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	var user User

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return User{}, fmt.Errorf("parse body: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, updUser UpdateUser) (User, error) {
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	l.Debug("update user...")
	err := s.repo.UpdateUser(ctx, updUser)
	if err != nil {
		return User{}, fmt.Errorf("update user: %w", err)
	}

	l.Debug("get user by ID...")
	return s.repo.UserByID(ctx, updUser.ID)
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	l.Debug("delete user...")
	return s.repo.DeleteUser(ctx, id, time.Now())
}

func (s *Service) StartWork(ctx context.Context, userID, taskID uuid.UUID) error {
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	_, err := s.repo.NotFinishedWorkHours(ctx, userID, taskID)
	if err == nil {
		return ErrWorkAlreadyStarted
	}

	if !errors.Is(err, ErrNotFound) {
		return err
	}

	wh := WorkHours{
		UserID:    userID,
		TaskID:    taskID,
		StartedAt: time.Now(),
	}

	l.Debug("start work...")
	return s.repo.StartWork(ctx, wh)
}

func (s *Service) FinishWork(ctx context.Context, userID, taskID uuid.UUID) error {
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	wh, err := s.repo.NotFinishedWorkHours(ctx, userID, taskID)
	if err != nil {
		return err
	}

	now := time.Now()
	wh.FinishedAt = &now
	wh.SpendTimeSec = int(wh.FinishedAt.Sub(wh.StartedAt).Seconds())

	l.Debug("finish work...")
	err = s.repo.FinishWork(ctx, wh)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) TaskSpendTimesByUser(ctx context.Context, id uuid.UUID, period Period) ([]TaskSpendTime, error) {
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	l.Debug("get task spend times by user...")
	spendTimesByUser, err := s.repo.TaskSpendTimesByUser(ctx, id, period)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return spendTimesByUser, nil
}

func (s *Service) Users(ctx context.Context, page, perPage int, filter UserFilter) ([]User, error) {
	l := ctx.Value(LoggerCtxKey{}).(*slog.Logger)

	l.Debug("get users...")
	return s.repo.Users(ctx, page, perPage, filter)
}
