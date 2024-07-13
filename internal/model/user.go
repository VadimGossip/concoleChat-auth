package model

import "time"

// Role enum
const (
	UnknownRole string = "unknown"
	UserRole    string = "user"
	AdminRole   string = "admin"
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
