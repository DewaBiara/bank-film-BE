package service

import (
	"context"

	"github.com/Budhiarta/bank-film-BE/internal/sharing/dto"
)

type SharingService interface {
	CreateSharing(ctx context.Context, sharing *dto.CreateSharing) error
	AddMember(ctx context.Context, userId string, sharing *dto.AddMemberRequest) error
	ValidateMember(ctx context.Context, sharing *dto.ValidateMember, userId string) (string, string, error)
}
