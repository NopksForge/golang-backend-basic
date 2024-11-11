package persistence

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserId    uuid.UUID
	UserEmail string
	UserName  string
	CreatedBy string
	CreatedAt time.Time
	UpdatedBy *string
	UpdatedAt *time.Time
}