package model

import (
	"database/sql"
	"time"
)

type UserInfo struct {
	Name     string `db:"username"`
	Email    string `db:"email"`
	Password string
	Role     string `db:"role"`
}

type User struct {
	ID        int64
	Info      UserInfo     `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
