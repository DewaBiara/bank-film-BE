package bootsrapper

import (
	"time"

	movieControllerPkg "github.com/Budhiarta/bank-film-BE/internal/movie/controller"
	movieRepositoryPkg "github.com/Budhiarta/bank-film-BE/internal/movie/repository/impl"
	movieServicePkg "github.com/Budhiarta/bank-film-BE/internal/movie/service/impl"
	sharingControllerPkg "github.com/Budhiarta/bank-film-BE/internal/sharing/controller"
	sharingRepositoryPkg "github.com/Budhiarta/bank-film-BE/internal/sharing/repository/impl"
	sharingServicePkg "github.com/Budhiarta/bank-film-BE/internal/sharing/service/impl"
	userControllerPkg "github.com/Budhiarta/bank-film-BE/internal/user/controller"
	userRepositoryPkg "github.com/Budhiarta/bank-film-BE/internal/user/repository/impl"
	userServicePkg "github.com/Budhiarta/bank-film-BE/internal/user/service/impl"
	"github.com/Budhiarta/bank-film-BE/pkg/routes"
	renderServicePkg "github.com/Budhiarta/bank-film-BE/pkg/utils/html/impl"
	jwtPkg "github.com/Budhiarta/bank-film-BE/pkg/utils/jwt_service/impl"
	passwordPkg "github.com/Budhiarta/bank-film-BE/pkg/utils/password/impl"
	qrPkg "github.com/Budhiarta/bank-film-BE/pkg/utils/qr/impl"
	smtp "github.com/Budhiarta/bank-film-BE/pkg/utils/smtp"
	"github.com/labstack/echo/v4"

	"gorm.io/gorm"
)

func InitController(e *echo.Echo, db *gorm.DB, conf map[string]string, mailer smtp.IMailer) {
	qrCodeService := qrPkg.NewCodeServiceImpl()
	renderService := renderServicePkg.NewRenderServiceImpl()
	passwordFunc := passwordPkg.NewPasswordFuncImpl()
	jwtService := jwtPkg.NewJWTService(conf["JWT_SECRET"], 1*time.Hour)

	// User
	userRepository := userRepositoryPkg.NewUserRepositoryImpl(db)
	userService := userServicePkg.NewUserServiceImpl(userRepository, passwordFunc, jwtService)
	userController := userControllerPkg.NewUserController(userService, jwtService)

	// OTP
	sharingRepository := sharingRepositoryPkg.NewSharingRepositoryImpl(db)
	sharingService := sharingServicePkg.NewSharingServiceImpl(sharingRepository, qrCodeService, renderService, jwtService, mailer, conf)
	sharingController := sharingControllerPkg.NewSharingController(sharingService, jwtService)

	// movie
	movieRepository := movieRepositoryPkg.NewMovieRepositoryImpl(db)
	movieService := movieServicePkg.NewMovieServiceImpl(movieRepository, jwtService)
	movieController := movieControllerPkg.NewMovieController(movieService, jwtService)

	route := routes.NewRoutes(userController, sharingController, movieController)
	route.Init(e, conf)
}
