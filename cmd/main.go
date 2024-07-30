package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"time-tracker/internal/app"
	"time-tracker/internal/tracker"
	"time-tracker/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx := context.Background()
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	cfg, err := app.NewConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	l.Info("ping Postgres...")
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	l.Info("ping Postgres OK")

	l.Info("up migrations...")
	err = upMigrations(cfg.PostgresDSN)
	if err != nil {
		log.Panic(err)
	}
	l.Info("up migrations OK")

	repo := tracker.NewRepository(db)
	service := tracker.NewService(repo, cfg.APIURL)
	handler := tracker.NewHandler(service)

	mw := tracker.NewMiddleware(l)

	router := http.NewServeMux()

	router.HandleFunc("POST /users", handler.CreateUser)
	router.HandleFunc("GET /users", handler.Users)
	router.HandleFunc("PATCH /users", handler.UpdateUser)
	router.HandleFunc("DELETE /users/{user_id}", handler.DeleteUser)
	router.HandleFunc("GET /users/{user_id}/report", handler.TaskSpendTimesByUser)

	router.HandleFunc("POST /work/start", handler.StartWork)
	router.HandleFunc("POST /work/finish", handler.FinishWork)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mw.Log(router),
		ReadTimeout:       time.Second * 3,
		ReadHeaderTimeout: time.Second,
	}

	go func() {
		err = server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}()

	l.Info("server started!", "port", cfg.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL)
	<-c

	err = server.Shutdown(ctx)
	if err != nil {
		log.Println("shutdown http server:", err)
	}
}

func upMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	fs := migrations.FS
	goose.SetBaseFS(fs)
	goose.SetLogger(goose.NopLogger())

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, ".")
	if err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		return err
	}

	return nil
}
