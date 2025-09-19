package post

import (
	"news-feed/internal/middleware"
)

type Service interface {
	Create(userID int, content string) (*Post, error)
	GetFeed(following []int, cursor int, limit int) ([]Post, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(userID int, content string) (*Post, error) {
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

func (s *service) GetFeed(userIDs []int, page, limit int) ([]Post, error) {
	var posts []Post
	offset := (page - 1) * limit
	err := s.repo.GetFeed(userIDs, offset, limit, &posts)
	return posts, err
}
