package impl

import (
	"crypto/rand"
	"crypto/rsa"

	// "crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	// "os"
)

// fungsi untuk membuat kunci publik dan kunci privat
func GenerateKeyPair() (privateKeyStr, publicKeyStr string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privateKeyPEM := pem.EncodeToMemory(privateKeyBlock)

	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyPEM := pem.EncodeToMemory(publicKeyBlock)

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// fungsi enkripsi pesan menggunakan public key
func Encrypt(publicKeyStr string, message []byte) ([]byte, error) {
	publicKeyBlock, _ := pem.Decode([]byte(publicKeyStr))
	if publicKeyBlock == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("parsed public key is not an RSA public key")
	}

	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, message)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message: %w", err)
	}

	modulusN := rsaPublicKey.N
	fmt.Println("modulus N: ", modulusN)

	return encryptedMessage, nil
}

// fungsi dekripsi pesan menggunakan private
func Decrypt(privateKeyStr string, encryptedMessage []byte) ([]byte, error) {
	privateKeyBlock, _ := pem.Decode([]byte(privateKeyStr))
	if privateKeyBlock == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded private key: %w", err)
	}

	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message: %w", err)
	}

	return decryptedMessage, nil
}
