package routes

import (
	// "fmt"
	// "net/http"

	movieControllerPkg "github.com/Budhiarta/bank-film-BE/internal/movie/controller"
	sharingControllerPkg "github.com/Budhiarta/bank-film-BE/internal/sharing/controller"
	userControllerPkg "github.com/Budhiarta/bank-film-BE/internal/user/controller"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/validation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Routes struct {
	userController    *userControllerPkg.UserController
	sharingController *sharingControllerPkg.SharingController
	movieController   *movieControllerPkg.MovieController
}

func NewRoutes(userController *userControllerPkg.UserController, sharingController *sharingControllerPkg.SharingController, movieController *movieControllerPkg.MovieController) *Routes {
	return &Routes{
		userController:    userController,
		sharingController: sharingController,
		movieController:   movieController,
	}
}

func (r *Routes) Init(e *echo.Echo, conf map[string]string) {
	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Validator = &validation.CustomValidator{Validator: validator.New()}

	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(conf["JWT_SECRET"]),
	})

	v1 := e.Group("/v1")

	// Users
	users := v1.Group("/users")
	users.POST("/signup/", r.userController.SignUpUser)
	users.POST("/login/", r.userController.LoginUser)

	usersWithAuth := users.Group("", jwtMiddleware)
	usersWithAuth.GET("/page/", r.userController.GetBriefUsers)
	usersWithAuth.PUT("/", r.userController.UpdateUser)
	usersWithAuth.GET("/detail/", r.userController.GetSingleUser)

	// Otp
	members := v1.Group("/members")

	membersWithAuth := members.Group("", jwtMiddleware)
	membersWithAuth.POST("/addmember/", r.sharingController.AddMember)
	membersWithAuth.POST("/validate/", r.sharingController.ValidateMember)

	// Movie
	movies := v1.Group("/movies")

	moviesWithAuth := movies.Group("", jwtMiddleware)
	moviesWithAuth.POST("/", r.movieController.CreateMovie, jwtMiddleware)
	moviesWithAuth.POST("/list/", r.movieController.AddListMovie, jwtMiddleware)
	moviesWithAuth.PUT("/", r.movieController.UpdateMovie, jwtMiddleware)
	moviesWithAuth.GET("/:movie_id/", r.movieController.GetSingleMovie, jwtMiddleware)
	moviesWithAuth.GET("/", r.movieController.GetPageMovie)
	moviesWithAuth.GET("/list/", r.movieController.GetListMovie)
	moviesWithAuth.DELETE("/:movie_id/", r.movieController.DeleteMovie, jwtMiddleware)
}
