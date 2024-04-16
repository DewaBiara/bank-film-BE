package impl

import (
	"context"

	"github.com/Budhiarta/bank-film-BE/internal/user/dto"
	"github.com/Budhiarta/bank-film-BE/internal/user/repository"
	"github.com/Budhiarta/bank-film-BE/internal/user/service"
	"github.com/Budhiarta/bank-film-BE/pkg/utils"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/jwt_service"
	otp "github.com/Budhiarta/bank-film-BE/pkg/utils/otp/impl"
	"github.com/Budhiarta/bank-film-BE/pkg/utils/password"
	"github.com/google/uuid"
)

type (
	UserServiceImpl struct {
		userRepository repository.UserRepository
		passwordHash   password.PasswordFunc
		jwtService     jwt_service.JWTService
	}
)

func NewUserServiceImpl(userRepository repository.UserRepository, function password.PasswordFunc, jwt jwt_service.JWTService) service.UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		passwordHash:   function,
		jwtService:     jwt,
	}
}

func (u *UserServiceImpl) SignUpUser(ctx context.Context, user *dto.UserSignUpRequest) error {
	hashedPassword, err := u.passwordHash.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	userEntity := user.ToEntity()
	userEntity.ID = uuid.New().String()

	err = u.userRepository.CreateUser(ctx, userEntity)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) LogInUser(ctx context.Context, user *dto.UserLoginRequest) (string, string, error) {
	userEntity, err := u.userRepository.FindByUsername(ctx, user.Username)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return "", "", utils.ErrInvalidCredentials
		}

		return "", "", err
	}

	err = u.passwordHash.CompareHashAndPassword([]byte(userEntity.Password), []byte(user.Password))
	if err != nil {
		return "", "", utils.ErrInvalidCredentials
	}

	token, err := u.jwtService.GenerateToken(userEntity)
	if err != nil {
		return "", "", err
	}

	otpCode := otp.GenerateRandomOTP()
	// fmt.Printf("Data: %s\n", otpCode)

	return token, otpCode, err
}

func (u *UserServiceImpl) GetSingleUser(ctx context.Context, userID string) (*dto.BriefUserResponse, error) {
	user, err := u.userRepository.GetSingleUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var userResponse = dto.NewBriefUserResponse(user)

	return userResponse, nil
}

func (u *UserServiceImpl) GetBriefUsers(ctx context.Context, page int, limit int) (*dto.BriefUsersResponse, error) {
	offset := (page - 1) * limit

	users, err := u.userRepository.GetBriefUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return dto.NewBriefUsersResponse(users), nil
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, userID string, request *dto.UserUpdateRequest) error {
	user := request.ToEntity()
	user.ID = userID

	if user.Password != "" {
		hashedPassword, err := u.passwordHash.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return err
		}

		user.Password = string(hashedPassword)
	}

	return u.userRepository.UpdateUser(ctx, user)
}
