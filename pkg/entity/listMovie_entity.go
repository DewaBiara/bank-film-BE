package entity

import (
	"time"

	"gorm.io/gorm"
)

type ListMovie struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	UserID    string
	Movies    *Movies `gorm:"many2many:Movie_ListMovie;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type ListMovies []ListMovie
