package entity

import (
	"time"

	"gorm.io/gorm"
)

type Sharing struct {
	ID            string `gorm:"primaryKey;type:varchar(36)"`
	SenderID      string `gorm:"type:varchar(255);not null;"`
	ReceiverID    string `gorm:"type:varchar(255);not null;"`
	ReceiverEmail string
	Chipertext    string
	PrivateKey    string
	PublicKey     string
	Otp           string
	OtpExpAt      time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Sharings []Sharing
