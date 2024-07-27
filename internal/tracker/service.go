package tracker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	user, err := s.getUserInfo(ctx, passportSeries, passportNumber)
	if err != nil {
		return fmt.Errorf("get user info: %w", err)
	}

	user.ID = uuid.Must(uuid.NewV4())
	user.PassportSeries = passportSeries
	user.PassportNumber = passportNumber
	user.CreatedAt = time.Now()

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
	err := s.repo.UpdateUser(ctx, updUser)
	if err != nil {
		return User{}, fmt.Errorf("update user: %w", err)
	}

	return s.repo.UserByID(ctx, updUser.ID)
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteUser(ctx, id, time.Now())
}

func (s *Service) StartWork(ctx context.Context, wh WorkHours) error {
	_, err := s.repo.NotFinishedWorkHours(ctx, wh.UserID, wh.TaskID)
	if err == nil {
		return ErrWorkAlreadyStarted
	}

	if !errors.Is(err, ErrNotFound) {
		return err
	}

	wh.StartedAt = time.Now()
	return s.repo.StartWork(ctx, wh)
}

func (s *Service) FinishWork(ctx context.Context, wh WorkHours) error {
	wh, err := s.repo.NotFinishedWorkHours(ctx, wh.UserID, wh.TaskID)
	if err != nil {
		return err
	}

	now := time.Now()
	wh.FinishedAt = &now
	wh.SpendTimeSec = int(wh.FinishedAt.Sub(wh.StartedAt).Seconds())

	err = s.repo.FinishWork(ctx, wh)
	if err != nil {
		return err
	}

	return nil
}
