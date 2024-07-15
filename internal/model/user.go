package model

import "time"

// Role
const (
	UnknownRole string = "UNKNOWN"
	UserRole    string = "USER"
	AdminRole   string = "ADMIN"
)

type UserInfo struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
}

type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UpdateUserInfo struct {
	Name  *string
	Email *string
	Role  string
}
