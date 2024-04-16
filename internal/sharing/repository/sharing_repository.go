package repository

import (
	"context"

	"github.com/Budhiarta/bank-film-BE/pkg/entity"
)

type SharingRepository interface {
	CreateSharing(ctx context.Context, sharing *entity.Sharing) error
	FindbyRecieverID(ctx context.Context, receiverID string, Otp string) (*entity.Sharing, error)
	AddMember(ctx context.Context, sharing *entity.Sharing) error
}
