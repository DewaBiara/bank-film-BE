package impl

import (
	"context"
	"strings"

	"github.com/Budhiarta/bank-film-BE/internal/user/repository"
	"github.com/Budhiarta/bank-film-BE/pkg/config"
	"github.com/Budhiarta/bank-film-BE/pkg/entity"
	"github.com/Budhiarta/bank-film-BE/pkg/utils"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) repository.UserRepository {
	userRepository := &UserRepositoryImpl{
		db: db,
	}

	err := userRepository.InitDefaultUser()
	if err != nil {
		panic(err)
	}

	return userRepository
}

func (u *UserRepositoryImpl) InitDefaultUser() error {
	var count int64
	err := u.db.Model(&entity.User{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count != 0 {
		return nil
	}

	err = u.db.Create(config.DefaultUser).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) CreateUser(ctx context.Context, user *entity.User) error {
	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			switch {
			case strings.Contains(err.Error(), "username"):
				return utils.ErrUsernameAlreadyExist
			}
		}

		return err
	}

	return nil
}

func (u *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).Select([]string{"id", "username", "password", "email", "age"}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *UserRepositoryImpl) GetSingleUser(ctx context.Context, userID string) (*entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).
		Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrMovieNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *UserRepositoryImpl) GetBriefUsers(ctx context.Context, limit int, offset int) (*entity.Users, error) {
	var users entity.Users
	err := u.db.WithContext(ctx).
		Select([]string{"id", "username", "name", "email", "age"}).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, utils.ErrUserNotFound
	}

	return &users, nil
}

func (u *UserRepositoryImpl) UpdateUser(ctx context.Context, user *entity.User) error {
	result := u.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		errStr := result.Error.Error()
		if strings.Contains(errStr, "Error 1062: Duplicate entry") {
			switch {
			case strings.Contains(errStr, "username"):
				return utils.ErrUsernameAlreadyExist
			}
		}

		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}
