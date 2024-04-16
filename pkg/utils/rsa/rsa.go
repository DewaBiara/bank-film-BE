package rsa

type rsaService interface {
	GenerateKeyPair() (string, string, error)
	Encrypt() (string, error)
	Decrypt() (string, error)
}
