package impl

import (
	"context"

	"github.com/Budhiarta/bank-film-BE/internal/movie/dto"
	"github.com/Budhiarta/bank-film-BE/internal/movie/repository"
	"github.com/Budhiarta/bank-film-BE/internal/movie/service"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/jwt_service"
	"github.com/google/uuid"
)

type (
	MovieServiceImpl struct {
		movieRepository repository.MovieRepository
		jwtService      jwt_service.JWTService
	}
)

func NewMovieServiceImpl(movieRepository repository.MovieRepository, jwtService jwt_service.JWTService) service.MovieService {
	return &MovieServiceImpl{
		movieRepository: movieRepository,
		jwtService:      jwtService,
	}
}

func (m *MovieServiceImpl) CreateMovie(ctx context.Context, movie *dto.CreateMovie) error {
	movieEntity := movie.ToEntity()
	movieEntity.ID = uuid.New().String()

	err := m.movieRepository.CreateMovie(ctx, movieEntity)
	if err != nil {
		return err
	}

	return nil
}

func (m *MovieServiceImpl) AddListMovie(ctx context.Context, movie *dto.AddToListRequest, userID string) error {
	movieEntity := movie.ToEntity()
	movieEntity.ID = uuid.New().String()
	movieEntity.UserID = userID

	err := m.movieRepository.AddListMovie(ctx, movieEntity)
	if err != nil {
		return err
	}

	return nil
}

func (m *MovieServiceImpl) GetSingleMovie(ctx context.Context, movieID string) (*dto.GetSingleMovieResponse, error) {
	item, err := m.movieRepository.GetSingleMovie(ctx, movieID)
	if err != nil {
		return nil, err
	}

	var movieResponse = dto.NewGetSingleMovieResponse(item)

	return movieResponse, nil
}

func (u *MovieServiceImpl) GetPageMovie(ctx context.Context, page int, limit int) (*dto.GetPageMoviesResponse, error) {
	offset := (page - 1) * limit

	movies, err := u.movieRepository.GetPageMovie(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return dto.NewGetPageMoviesResponse(movies), nil
}

func (u *MovieServiceImpl) GetListMovie(ctx context.Context, page int, limit int, userID string, ageUser int64) (*dto.GetListMoviesResponse, error) {
	offset := (page - 1) * limit

	movies, err := u.movieRepository.GetListMovie(ctx, limit, offset, userID, ageUser)
	if err != nil {
		return nil, err
	}

	return dto.NewGetListMoviesResponse(movies), nil
}

func (u *MovieServiceImpl) UpdateMovie(ctx context.Context, movieID string, updateMovie *dto.UpdateMovieRequest) error {
	movie := updateMovie.ToEntity()
	movie.ID = movieID

	return u.movieRepository.UpdateMovie(ctx, movie)
}

func (d *MovieServiceImpl) DeleteMovie(ctx context.Context, movieID string) error {

	return d.movieRepository.DeleteMovie(ctx, movieID)
}
