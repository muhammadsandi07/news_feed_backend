package user

import (
	"news-feed/internal/middleware"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(username, password string) (*User, error)
	Login(username, password string) (string, string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Register(username, password string) (*User, error) {
	if _, err := s.repo.FindByUsername(username); err == nil {
		return nil, middleware.Conflict("username already exists")
	}
	if len(password) < 6 {
		return nil, middleware.UnprocessableEntity("password must be at least 6 chars")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	u := &User{
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *service) Login(username, password string) (string, string, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", "", middleware.NotFound("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", middleware.Unauthorized("invalid credentials")
	}

	accessToken, err := generateToken(u.ID, os.Getenv("JWT_SECRET"), time.Minute*15)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateToken(u.ID, os.Getenv("JWT_SECRET"), time.Hour*24*7)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateToken(userID string, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
