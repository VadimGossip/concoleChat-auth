package app

import (
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	"github.com/VadimGossip/concoleChat-auth/internal/repository/user"
)

type DBAdapter struct {
	userRepo repository.UserRepository
}

func NewDBAdapter() *DBAdapter {
	return &DBAdapter{userRepo: user.NewRepository()}
}
