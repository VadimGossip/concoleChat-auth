package model

import (
	"database/sql"
	"time"
)

type UserInfo struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
}

type User struct {
	Id        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
