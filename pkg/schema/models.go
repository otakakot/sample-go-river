// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package schema

import (
	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
