package post

import (
	"news-feed/internal/middleware"
)

type Service interface {
	Create(userID uint, content string) (*Post, error)
	GetFeed(following []uint, cursor string, limit int) ([]Post, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(userID uint, content string) (*Post, error) {
	if len(content) == 0 {
		return nil, middleware.BadRequest("content cannot be empty")
	}
	if len(content) > 200 {
		return nil, middleware.UnprocessableEntity("content must be at most 200 characters")
	}

	p := &Post{
		UserID:  userID,
		Content: content,
	}

	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) GetFeed(following []uint, cursor string, limit int) ([]Post, error) {
	if len(following) == 0 {
		return []Post{}, nil
	}
	return s.repo.GetFeed(following, cursor, limit)
}
