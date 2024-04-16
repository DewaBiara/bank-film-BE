package dto

import "github.com/Budhiarta/bank-film-BE/pkg/entity"

type UserSignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Age      int64  `json:"age"`
	Telp     string `json:"telp" validate:"required"`
}

func (u *UserSignUpRequest) ToEntity() *entity.User {
	return &entity.User{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Name:     u.Name,
		Age:      u.Age,
		Telp:     u.Telp,
	}
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      int64  `json:"age"`
	Telp     string `json:"telp"`
}

func (u *UserUpdateRequest) ToEntity() *entity.User {
	return &entity.User{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Name:     u.Name,
		Age:      u.Age,
		Telp:     u.Telp,
	}
}

type BriefUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      int64  `json:"age"`
	Telp     string `json:"telp"`
}

func NewBriefUserResponse(user *entity.User) *BriefUserResponse {
	return &BriefUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		Age:      user.Age,
		Telp:     user.Telp,
	}
}

type BriefUsersResponse []BriefUserResponse

func NewBriefUsersResponse(users *entity.Users) *BriefUsersResponse {
	var briefUsersResponse BriefUsersResponse
	for _, user := range *users {
		briefUsersResponse = append(briefUsersResponse, *NewBriefUserResponse(&user))
	}
	return &briefUsersResponse
}

type ValidateOTPReq struct {
	Token string `json:"token"`
	OTP   string `json:"otp"`
}
