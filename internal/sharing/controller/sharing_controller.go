package controller

import (
	"net/http"

	"github.com/Budhiarta/bank-film-BE/internal/sharing/dto"
	"github.com/Budhiarta/bank-film-BE/internal/sharing/service"
	"github.com/Budhiarta/bank-film-BE/pkg/utils"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/jwt_service"

	"github.com/labstack/echo/v4"
)

type SharingController struct {
	sharingService service.SharingService
	jwtService     jwt_service.JWTService
}

func NewSharingController(sharingService service.SharingService, jwtService jwt_service.JWTService) *SharingController {
	return &SharingController{
		sharingService: sharingService,
		jwtService:     jwtService,
	}
}

func (s *SharingController) AddMember(c echo.Context) error {
	claims := s.jwtService.GetClaims(&c)
	userID := claims["user_id"].(string)

	sharing := new(dto.AddMemberRequest)
	if err := c.Bind(sharing); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(sharing); err != nil {
		return err
	}

	err := s.sharingService.AddMember(c.Request().Context(), userID, sharing)
	if err != nil {
		switch err {
		case utils.ErrUsernameAlreadyExist:
			fallthrough
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "success Adding Member",
	})
}

func (s *SharingController) ValidateMember(c echo.Context) error {
	claims := s.jwtService.GetClaims(&c)
	userID := claims["user_id"].(string)

	sharing := new(dto.ValidateMember)
	if err := c.Bind(sharing); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(sharing); err != nil {
		return err
	}

	senderId, status, err := s.sharingService.ValidateMember(c.Request().Context(), sharing, userID)
	if err != nil {
		switch err {
		case utils.ErrOtpInvalid:
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		case utils.ErrOtpExpired:
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		// "message": "verification succes",
		"status":    status,
		"id_sender": senderId,
	})
}
