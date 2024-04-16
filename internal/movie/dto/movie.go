package dto

import "github.com/Budhiarta/bank-film-BE/pkg/entity"

type CreateMovie struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	AgeRestriction int64  `json:"age_restriction"`
}

func (m *CreateMovie) ToEntity() *entity.Movie {
	return &entity.Movie{
		Title:          m.Title,
		Description:    m.Description,
		AgeRestriction: m.AgeRestriction,
	}
}

type AddToListRequest struct {
	UserID string         `json:"userid"`
	Movies *MovieRequests `json:"movies"`
}

func (m *AddToListRequest) ToEntity() *entity.ListMovie {
	return &entity.ListMovie{
		UserID: m.UserID,
		Movies: m.Movies.ToEntity(),
	}
}

type UpdateMovieRequest struct {
	ID             string `json:"id" validate:"required"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AgeRestriction int64  `json:"age_restriction"`
}

func (m *UpdateMovieRequest) ToEntity() *entity.Movie {
	return &entity.Movie{
		Title:          m.Title,
		Description:    m.Description,
		AgeRestriction: m.AgeRestriction,
	}
}

type GetSingleMovieResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AgeRestriction int64  `json:"age_restriction"`
}

func NewGetSingleMovieResponse(movie *entity.Movie) *GetSingleMovieResponse {
	return &GetSingleMovieResponse{
		ID:             movie.ID,
		Title:          movie.Title,
		Description:    movie.Description,
		AgeRestriction: movie.AgeRestriction,
	}
}

type GetPageMovieResponse struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AgeRestriction int64  `json:"age_restriction"`
}

func NewGetPageMovieResponse(movie *entity.Movie) *GetPageMovieResponse {
	return &GetPageMovieResponse{
		ID:             movie.ID,
		Title:          movie.Title,
		Description:    movie.Description,
		AgeRestriction: movie.AgeRestriction,
	}
}

type GetPageMoviesResponse []GetPageMovieResponse

func NewGetPageMoviesResponse(movies *entity.Movies) *GetPageMoviesResponse {
	var getPageMovies GetPageMoviesResponse
	for _, movies := range *movies {
		getPageMovies = append(getPageMovies, *NewGetPageMovieResponse(&movies))
	}
	return &getPageMovies
}

type GetListMovieResponse struct {
	ID     string          `json:"id"`
	UserID string          `json:"userid"`
	Movies *MovieResponses `json:"movies"`
}

func NewGetListMovieResponse(movie *entity.ListMovie) *GetListMovieResponse {
	return &GetListMovieResponse{
		ID:     movie.ID,
		UserID: movie.UserID,
		Movies: NewMovieResponses(movie.Movies),
	}
}

type GetListMoviesResponse []GetListMovieResponse

func NewGetListMoviesResponse(movie *entity.ListMovies) *GetListMoviesResponse {
	var getListMovies GetListMoviesResponse
	for _, movie := range *movie {
		getListMovies = append(getListMovies, *NewGetListMovieResponse(&movie))
	}
	return &getListMovies
}

type MovieRequest struct {
	MovieID string `json:"movieid"`
}

type MovieRequests []MovieRequest

func (r MovieRequest) ToEntity() *entity.Movie {
	return &entity.Movie{
		ID: r.MovieID,
	}
}

func (r *MovieRequests) ToEntity() *entity.Movies {
	movies := entity.Movies{}
	for _, movie := range *r {
		movies = append(movies, *movie.ToEntity())
	}

	return &movies
}

type MovieResponse struct {
	MovieID        string `json:"movieid"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	AgeRestriction int64  `json:"age_restriction"`
}

func NewMovieResponse(movie *entity.Movie) *MovieResponse {
	return &MovieResponse{
		MovieID:        movie.ID,
		Title:          movie.Title,
		Description:    movie.Description,
		AgeRestriction: movie.AgeRestriction,
	}
}

type MovieResponses []MovieResponse

func NewMovieResponses(movies *entity.Movies) *MovieResponses {
	var movieResponses MovieResponses
	for _, movies := range *movies {
		movieResponses = append(movieResponses, *NewMovieResponse(&movies))
	}
	return &movieResponses
}
