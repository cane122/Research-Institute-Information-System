package services

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/cane/research-institute-system/backend/models"
	"github.com/cane/research-institute-system/backend/repositories"

	"golang.org/x/crypto/argon2"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User    *models.User `json:"user"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return &LoginResponse{
			Success: false,
			Message: "Korisničko ime i lozinka su obavezni",
		}, nil
	}

	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return &LoginResponse{
			Success: false,
			Message: "Neispravno korisničko ime ili lozinka",
		}, nil
	}

	if user.Status != "aktivan" {
		return &LoginResponse{
			Success: false,
			Message: "Nalog nije aktivan",
		}, nil
	}

	// Check if this is first-time login (no previous login recorded)
	if user.PoslednajaPrijava == nil {
		return &LoginResponse{
			User:    user,
			Success: true,
			Message: "FIRST_TIME_LOGIN", // Special message indicating first-time login
		}, nil
	}

	if !s.verifyPassword(req.Password, user.HashSifre) {
		return &LoginResponse{
			Success: false,
			Message: "Neispravno korisničko ime ili lozinka",
		}, nil
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.KorisnikID)

	// Clear password hash from response
	user.HashSifre = ""

	return &LoginResponse{
		User:    user,
		Success: true,
		Message: "Uspešna prijava",
	}, nil
}

func (s *AuthService) CreateUser(user *models.User, tempPassword string) error {
	if user.KorisnickoIme == "" || user.Email == "" {
		return errors.New("korisničko ime i email su obavezni")
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByUsername(user.KorisnickoIme)
	if existingUser != nil {
		return errors.New("korisnik sa tim korisničkim imenom već postoji")
	}

	// Hash the temporary password
	hashedPassword, err := s.HashPassword(tempPassword)
	if err != nil {
		return err
	}

	user.HashSifre = hashedPassword
	user.Status = "aktivan"

	return s.userRepo.Create(user)
}

func (s *AuthService) ResetPassword(userID int) (string, error) {
	// Generate temporary password
	tempPassword := s.generateTempPassword()

	hashedPassword, err := s.HashPassword(tempPassword)
	if err != nil {
		return "", err
	}

	err = s.userRepo.UpdatePassword(userID, hashedPassword)
	if err != nil {
		return "", err
	}

	return tempPassword, nil
}

func (s *AuthService) ChangePassword(userID int, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("lozinka mora imati najmanje 8 karaktera")
	}

	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

func (s *AuthService) CompleteFirstTimeSetup(userID int, newPassword string) error {
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password and set first login timestamp
	err = s.userRepo.UpdatePassword(userID, hashedPassword)
	if err != nil {
		return err
	}

	// Mark as having completed first login
	return s.userRepo.UpdateLastLogin(userID)
}

func (s *AuthService) CompleteFirstTimeSetupByUsername(username, newPassword string) error {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return errors.New("korisnik nije pronađen")
	}

	// Use existing function with userID
	return s.CompleteFirstTimeSetup(user.KorisnikID, newPassword)
}

func (s *AuthService) HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, 64*1024, 1, 4, b64Salt, b64Hash)
	return full, nil
}

func (s *AuthService) verifyPassword(password, hash string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false
	}
	if version != argon2.Version {
		return false
	}

	var memory, iterations, parallelism int
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, uint32(iterations), uint32(memory), uint8(parallelism), uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1
}

func (s *AuthService) generateTempPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	rand.Read(b)
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
