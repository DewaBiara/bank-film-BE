package repository

import (
	"context"

	"github.com/Budhiarta/bank-film-BE/pkg/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	GetSingleUser(ctx context.Context, userID string) (*entity.User, error)
	GetBriefUsers(ctx context.Context, limit int, offset int) (*entity.Users, error)
	UpdateUser(ctx context.Context, user *entity.User) error
}
