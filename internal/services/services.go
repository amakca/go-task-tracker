package service

import (
	"go-task-tracker/internal/repo"
	"go-task-tracker/internal/services/contracts"
	"go-task-tracker/internal/services/users"
	"go-task-tracker/pkg/hasher"
	"time"
)

type Services struct {
	Auth contracts.Auth
}

type ServicesDependencies struct {
	Repos *repo.Repositories
	// GDrive webapi.GDrive
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth: users.NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
	}
}
