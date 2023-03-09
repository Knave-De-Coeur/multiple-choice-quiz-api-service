package api

import (
	"database/sql"
)

type User struct {
	ID                 string          `json:"ID"`
	Name               string          `json:"name" validate:"required"`
	Age                int8            `json:"age" validate:"required"`
	Username           string          `json:"username" validate:"required"`
	Password           string          `json:"password,omitempty" validate:"required"`
	CreatedAT          *sql.NullString `json:"created_at,omitempty"`
	UpdatedAT          *sql.NullString `json:"updated_at,omitempty"`
	LastLoginTimeStamp *sql.NullTime   `json:"last_login_time_stamp,omitempty"`
}
