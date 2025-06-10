package services

import (
	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/repositories"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	LoginUser(input *models.LoginRequest) (string, string, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) LoginUser(input *models.LoginRequest) (string, string, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return "", "", err
	}
	
	// Compare the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", "", err // Password mismatch
	}
	
	// Generate JWT token
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", err // Error generating token
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err // Error generating refresh token
	}
	return accessToken, refreshToken, nil
}