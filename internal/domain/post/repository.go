package post

import "gorm.io/gorm"

type Repository interface {
	Create(post *Post) error
	GetFeed(userIDs []uint, cursor string, limit int) ([]Post, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(post *Post) error {
	return r.db.Create(post).Error
}

func (r *repository) GetFeed(userIDs []uint, cursor string, limit int) ([]Post, error) {
	var posts []Post

	q := r.db.Where("user_id IN ?", userIDs)

	if cursor != "" {
		q = q.Where("created_at < ?", cursor)
	}

	err := q.Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error

	return posts, err
}
