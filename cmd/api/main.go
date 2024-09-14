package main

import (
	"cmp"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"github.com/otakakot/sample-go-river/internal/riverx"
	"github.com/otakakot/sample-go-river/pkg/schema"
)

func main() {
	conn := cmp.Or(os.Getenv("DSN"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	pool, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	riverCli, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		panic(err)
	}

	port := cmp.Or(os.Getenv("PORT"), "8080")

	hdl := http.NewServeMux()

	hdl.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ユーザーからの入力をごにょごにょする実装はサボりました

		ctx := r.Context()

		tx, err := pool.Begin(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
		defer tx.Rollback(ctx)

		email := uuid.New().String() + "@example.com"

		us, err := schema.New(tx).InsertUser(ctx, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		res, err := riverCli.InsertTx(
			ctx,
			tx,
			riverx.JobArgs{
				UserID: us.ID,
			},
			&river.InsertOpts{},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		if err := tx.Commit(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusAccepted)

		w.Write([]byte(strconv.Itoa(int(res.Job.ID))))
	})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           hdl,
		ReadHeaderTimeout: 30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("start server listen")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	slog.Info("start server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	slog.Info("done server shutdown")
}
