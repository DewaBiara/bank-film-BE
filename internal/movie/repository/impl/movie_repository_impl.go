package impl

import (
	"context"
	"errors"
	"strings"

	"github.com/Budhiarta/bank-film-BE/internal/movie/repository"
	"github.com/Budhiarta/bank-film-BE/pkg/entity"
	"github.com/Budhiarta/bank-film-BE/pkg/utils"
	"gorm.io/gorm"
)

type MovieRepositoryImpl struct {
	db *gorm.DB
}

func NewMovieRepositoryImpl(db *gorm.DB) repository.MovieRepository {
	movieRepository := &MovieRepositoryImpl{
		db: db,
	}

	return movieRepository
}

func (m *MovieRepositoryImpl) CreateMovie(ctx context.Context, movie *entity.Movie) error {
	err := m.db.WithContext(ctx).Create(movie).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *MovieRepositoryImpl) AddListMovie(ctx context.Context, movie *entity.ListMovie) error {
	err := m.db.WithContext(ctx).Create(movie).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062: Duplicate entry") {
			switch {
			case strings.Contains(err.Error(), "name"):
				return err
			}
		}

		return err
	}
	return nil
}

// func (m *MovieRepositoryImpl) AddToListRequest(ctx context.Context, title string)  error {
// 	err := m.db.WithContext(ctx).Create(movie).Error
// 	if err != nil{
// 		return err
// 	}

// 	return nil
// }

func (m *MovieRepositoryImpl) UpdateMovie(ctx context.Context, movie *entity.Movie) error {
	result := m.db.WithContext(ctx).Model(&entity.Movie{}).Where("id = ?", movie.ID).Updates(movie)
	if result.Error != nil {
		errStr := result.Error.Error()
		if strings.Contains(errStr, "Error 1062: Duplicate entry") {
			switch {
			case strings.Contains(errStr, "title"):
				return utils.ErrTitleAlreadyExist
			}
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return utils.ErrMovieNotFound
	}
	return nil
}

func (m *MovieRepositoryImpl) GetSingleMovie(ctx context.Context, movieID string) (*entity.Movie, error) {
	var movie entity.Movie
	err := m.db.WithContext(ctx).Select([]string{"id", "title", "description", "age_restriction"}).
		Where("id = ?", movieID).First(&movie).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrMovieNotFound
		}

		return nil, err
	}

	return &movie, nil
}

func (m *MovieRepositoryImpl) GetPageMovie(ctx context.Context, limit int, offset int) (*entity.Movies, error) {
	var movies entity.Movies
	err := m.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&movies).Error
	if err != nil {
		return nil, err
	}
	if len(movies) == 0 {
		return nil, utils.ErrMovieNotFound
	}
	return &movies, nil
}

func (m *MovieRepositoryImpl) GetListMovie(ctx context.Context, limit int, offset int, userID string, ageUser int64) (*entity.ListMovies, error) {
	var movies entity.ListMovies
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Preload("Movies", "age_restriction <= ?", ageUser). // Adding WHERE condition to preload
		Find(&movies).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrMovieNotFound
		}
		return nil, err
	}
	if len(movies) == 0 {
		return nil, utils.ErrMovieNotFound
	}
	return &movies, nil
}

func (d *MovieRepositoryImpl) DeleteMovie(ctx context.Context, movieID string) error {
	result := d.db.WithContext(ctx).Select("Movie").Delete(&entity.Movie{}, "id = ?", movieID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrMovieNotFound
	}
	return nil
}
