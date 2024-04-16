package controller

import (
	"net/http"
	"strconv"

	"github.com/Budhiarta/bank-film-BE/internal/movie/dto"
	"github.com/Budhiarta/bank-film-BE/internal/movie/service"
	"github.com/Budhiarta/bank-film-BE/pkg/utils"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/jwt_service"
	"github.com/labstack/echo/v4"
)

type MovieController struct {
	movieService service.MovieService
	jwtService   jwt_service.JWTService
}

func NewMovieController(movieService service.MovieService, jwtService jwt_service.JWTService) *MovieController {
	return &MovieController{
		movieService: movieService,
		jwtService:   jwtService,
	}
}

func (m *MovieController) CreateMovie(c echo.Context) error {
	movie := new(dto.CreateMovie)
	if err := c.Bind(movie); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(movie); err != nil {
		return err
	}

	err := m.movieService.CreateMovie(c.Request().Context(), movie)
	if err != nil {
		switch err {
		case utils.ErrTitleAlreadyExist:
			fallthrough
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Success Creating Movie",
	})
}

func (m *MovieController) AddListMovie(c echo.Context) error {
	claims := m.jwtService.GetClaims(&c)
	userID := claims["user_id"].(string)

	movie := new(dto.AddToListRequest)
	if err := c.Bind(movie); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(movie); err != nil {
		return err
	}

	err := m.movieService.AddListMovie(c.Request().Context(), movie, userID)
	if err != nil {
		switch err {
		case utils.ErrTitleAlreadyExist:
			fallthrough
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Success Adding List Movie",
	})
}

func (u *MovieController) UpdateMovie(c echo.Context) error {
	// claims := u.jwtService.GetClaims(&c)
	movie := new(dto.UpdateMovieRequest)
	if err := c.Bind(movie); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(movie); err != nil {
		return err
	}

	err := u.movieService.UpdateMovie(c.Request().Context(), movie.ID, movie)
	if err != nil {
		switch err {
		case utils.ErrMovieNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case utils.ErrTitleAlreadyExist:
			fallthrough
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success update movie",
		"data":    movie,
	})
}

func (u *MovieController) GetSingleMovie(c echo.Context) error {
	movieID := c.Param("movie_id")
	movie, err := u.movieService.GetSingleMovie(c.Request().Context(), movieID)
	if err != nil {
		if err == utils.ErrMovieNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting movie",
		"data":    movie,
	})
}

func (u *MovieController) GetPageMovie(c echo.Context) error {

	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidNumber.Error())
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "20"
	}
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidNumber.Error())
	}

	movie, err := u.movieService.GetPageMovie(c.Request().Context(), int(pageInt), int(limitInt))
	if err != nil {
		if err == utils.ErrMovieNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting movies",
		"data":    movie,
		"meta": echo.Map{
			"page":  pageInt,
			"limit": limitInt,
		},
	})
}

func (u *MovieController) GetListMovie(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userID := c.QueryParam("user_id")
	if userID == "" {
		userId := claims["user_id"].(string)
		userID = userId
	}

	var ageUser int64
	if ageClaim, ok := claims["age"].(float64); ok {
		ageUser = int64(ageClaim)
	} else {
		// Jika klaim umur tidak ada atau nil, kembalikan respons dengan kesalahan yang sesuai
		return echo.NewHTTPError(http.StatusBadRequest, "Age claim not found or invalid")
	}

	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidNumber.Error())
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "20"
	}
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidNumber.Error())
	}

	movie, err := u.movieService.GetListMovie(c.Request().Context(), int(pageInt), int(limitInt), userID, ageUser)
	if err != nil {
		if err == utils.ErrMovieNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success getting movies",
		"data":    movie,
		"meta": echo.Map{
			"page":  pageInt,
			"limit": limitInt,
		},
	})
}

func (d *MovieController) DeleteMovie(c echo.Context) error {
	// claims := d.jwtService.GetClaims(&c)
	// // role := claims["role"].(string)
	// // if role != "admin" {
	// // 	return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	// // }
	movieID := c.Param("movie_id")
	err := d.movieService.DeleteMovie(c.Request().Context(), movieID)
	if err != nil {
		switch err {
		case utils.ErrMovieNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success deleting movie",
	})
}
