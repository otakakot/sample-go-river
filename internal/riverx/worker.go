package riverx

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"

	"github.com/otakakot/sample-go-river/pkg/schema"
)

var _ river.Worker[JobArgs] = (*Worker)(nil)

type Worker struct {
	pool *pgxpool.Pool
	river.WorkerDefaults[JobArgs]
}

func New(pool *pgxpool.Pool) *Worker {
	return &Worker{
		pool: pool,
	}
}

// Work implements river.Worker.
func (w *Worker) Work(ctx context.Context, job *river.Job[JobArgs]) error {
	slog.Info("start work" + strconv.Itoa(int(job.ID)))
	defer slog.Info("end work" + strconv.Itoa(int(job.ID)))

	us, err := schema.New(w.pool).FindUserByID(ctx, job.Args.UserID)
	if err != nil {
		return fmt.Errorf("find user by id: %w", err)
	}

	slog.Info("start send email: " + us.Email)

	time.Sleep(5 * time.Second)

	slog.Info("done send email: " + us.Email)

	return nil
}
