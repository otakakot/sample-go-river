package riverx

import (
	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

var _ river.JobArgs = (*JobArgs)(nil)

type JobArgs struct {
	UserID uuid.UUID
}

// Kind implements river.JobArgs.
func (ja JobArgs) Kind() string {
	return "job"
}
