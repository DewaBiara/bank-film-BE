package service

import (
	"context"

	"github.com/Budhiarta/bank-film-BE/internal/movie/dto"
)

type MovieService interface {
	CreateMovie(ctx context.Context, movie *dto.CreateMovie) error
	AddListMovie(ctx context.Context, movie *dto.AddToListRequest, userID string) error
	UpdateMovie(ctx context.Context, movieID string, updateMovie *dto.UpdateMovieRequest) error
	GetSingleMovie(ctx context.Context, movieID string) (*dto.GetSingleMovieResponse, error)
	GetPageMovie(ctx context.Context, limit int, offset int) (*dto.GetPageMoviesResponse, error)
	GetListMovie(ctx context.Context, limit int, offset int, userID string, ageUser int64) (*dto.GetListMoviesResponse, error)
	DeleteMovie(ctx context.Context, movieID string) error
}
