package entity

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID             string `gorm:"primaryKey;type:varchar(36)"`
	Title          string `gorm:"type:varchar(255);not null;"`
	Description    string `gorm:"type:varchar(255);not null;"`
	AgeRestriction int64
	ListMovies     ListMovies `gorm:"many2many:Movie_ListMovie;"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Movies []Movie
