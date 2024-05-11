package infrastructure

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTPMultiFactor struct {
	cfg            *Config
	userRepository domain.UserRepository
}

func NewTOTPMultiFactor(cfg *Config, userRepo domain.UserRepository) *TOTPMultiFactor {
	return &TOTPMultiFactor{
		cfg:            cfg,
		userRepository: userRepo,
	}
}

// Documentation: https://github.com/pquerna/otp/blob/master/doc.go
// GenerateTOTP generates a time-based OTP for multi factor authentication
func (m TOTPMultiFactor) GenerateTOTP(ctx context.Context, userID int, email string) (*otp.Key, string, error) {
	totpOptions := totp.GenerateOpts{
		Issuer:      m.cfg.TOTP.Issuer,
		AccountName: email,
		SecretSize:  10,
	}
	key, err := totp.Generate(totpOptions)
	if err != nil {
		log.Printf("failed to generate totp for user with email %s with error: %v\n", email, err)
		return nil, "", err
	}

	// Encrypt the TOTP secret
	encryptedSecret, err := m.encryptTOTPSecret(key.Secret(), []byte(m.cfg.TOTP.EncryptionKey))
	if err != nil {
		return nil, "", err
	}

	totpConfiguration := domain.TOTPConfiguration{
		UserID:              userID,
		Email:               email,
		TOTPEncryptedSecret: encryptedSecret,
		CreatedAt:           time.Now(),
	}

	if err := m.userRepository.InsertUserTOTPSecret(ctx, totpConfiguration); err != nil {
		log.Println("failed to insert user totp secret with error:", err)
		return nil, "", err
	}

	return key, encryptedSecret, nil
}

func (m TOTPMultiFactor) encryptTOTPSecret(secret string, encryptionKey []byte) (string, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		log.Println("failed to create AES cipher block with error:", err)
		return "", err
	}

	// Pad the secret if necessary to match block size
	plainText := []byte(secret)
	paddedPlainText := make([]byte, len(plainText))
	copy(paddedPlainText, plainText)

	// Create a new AES GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("failed to create new AES GCM cipher with error:", err)
		return "", err
	}

	// Generate a nonce (to prevent replay attacks)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("failed to perform read nonce value with error:", err)
		return "", err
	}

	// encrypt the secret with encryption key
	cipherText := gcm.Seal(nonce, nonce, paddedPlainText, nil)

	// Encode the cipher text to base64
	encryptedSecret := base64.StdEncoding.EncodeToString(cipherText)

	return encryptedSecret, nil
}

func (m TOTPMultiFactor) decryptTOTPSecret(encryptedSecret string, encryptionKey []byte) (string, error) {
	return "", nil
}
