package repository

import (
	"context"

	"github.com/Budhiarta/bank-film-BE/pkg/entity"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *entity.Movie) error
	AddListMovie(ctx context.Context, movie *entity.ListMovie) error
	// AddToListRequest(ctx context.Context, title string) error
	UpdateMovie(ctx context.Context, movie *entity.Movie) error
	GetSingleMovie(ctx context.Context, movieID string) (*entity.Movie, error)
	GetPageMovie(ctx context.Context, limit int, offset int) (*entity.Movies, error)
	GetListMovie(ctx context.Context, limit int, offset int, userID string, ageUser int64) (*entity.ListMovies, error)
	DeleteMovie(ctx context.Context, movieID string) error
}
