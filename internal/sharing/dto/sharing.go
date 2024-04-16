package dto

import "github.com/Budhiarta/bank-film-BE/pkg/entity"

type CreateSharing struct {
	SenderID   string `json:"sender_id"`
	Otp        string `json:"otp"`
	Chipertext string `json:"chipertext"`
	PrivateKey string `json:"privatekey"`
	PublicKey  string `json:"publickey"`
}

func (s *CreateSharing) ToEntity() *entity.Sharing {
	return &entity.Sharing{
		SenderID:   s.SenderID,
		Otp:        s.Otp,
		Chipertext: s.Chipertext,
		PrivateKey: s.PrivateKey,
		PublicKey:  s.PublicKey,
	}
}

type AddMemberRequest struct {
	ReceiverID    string `json:"receiver_id" validate:"required"`
	ReceiverEmail string `json:"receiver_email"`
}

func (s *AddMemberRequest) ToEntity() *entity.Sharing {
	return &entity.Sharing{
		ReceiverID:    s.ReceiverID,
		ReceiverEmail: s.ReceiverEmail,
	}
}

type ValidateMember struct {
	ID         string `json:"id"`
	Chipertext string `json:"chipertext"`
}
