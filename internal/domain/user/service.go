package user

import (
	"fmt"
	"news-feed/internal/middleware"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(username, password string) (*User, error)
	Login(username, password string) (string, string, error)
	Refresh(refreshToken string) (string, error)
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{repo, jwtSecret}
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

	accessToken, err := generateToken(u.ID, s.jwtSecret, time.Minute*15)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateToken(u.ID, s.jwtSecret, time.Hour*24*7)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) Refresh(refreshToken string) (string, error) {
	claims, err := parseToken(refreshToken, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", middleware.Unauthorized("invalid refresh token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", middleware.Unauthorized("invalid refresh token claims")
	}

	// generate new access token
	accessToken, err := generateToken(string(userID), s.jwtSecret, time.Minute*15)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func generateToken(userID string, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func parseToken(tokenStr, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
