package main

import (
	"cmp"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"github.com/otakakot/sample-go-river/internal/riverx"
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

	worker := river.NewWorkers()

	river.AddWorker(worker, riverx.New(pool))

	cli, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {
				MaxWorkers: 1,
			},
		},
		Workers: worker,
	})
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("start worker")

		if err := cli.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()

	slog.Info("start worker shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := cli.Stop(ctx); err != nil {
		panic(err)
	}

	slog.Info("done worker shutdown")
}
