package otp

type OTPService interface {
	GenerateRandomOTP() (string, error)
}
