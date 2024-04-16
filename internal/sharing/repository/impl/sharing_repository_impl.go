package impl

import (
	"context"
	"strings"

	"github.com/Budhiarta/bank-film-BE/internal/sharing/repository"
	"github.com/Budhiarta/bank-film-BE/pkg/entity"
	"github.com/Budhiarta/bank-film-BE/pkg/utils"
	"gorm.io/gorm"
)

type SharingRepositoryImpl struct {
	db *gorm.DB
}

func NewSharingRepositoryImpl(db *gorm.DB) repository.SharingRepository {
	sharingRepository := &SharingRepositoryImpl{
		db: db,
	}

	return sharingRepository
}

func (s *SharingRepositoryImpl) CreateSharing(ctx context.Context, sharing *entity.Sharing) error {
	err := s.db.WithContext(ctx).Create(sharing).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *SharingRepositoryImpl) AddMember(ctx context.Context, sharing *entity.Sharing) error {
	err := s.db.WithContext(ctx).Create(sharing).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate Entry") {
			switch {
			case strings.Contains(err.Error(), "receiver_id"):
				return utils.ErrUsernameAlreadyExist
			}
		}
		return err
	}

	return nil
}

func (s *SharingRepositoryImpl) FindbyRecieverID(ctx context.Context, receiverID string, id string) (*entity.Sharing, error) {
	var sharing entity.Sharing
	err := s.db.WithContext(ctx).Where("receiver_id = ?", receiverID).Where("id = ?", id).First(&sharing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFound
		}

		return nil, err
	}

	return &sharing, nil
}
