package post

import "gorm.io/gorm"

type Repository interface {
	Create(post *Post) error
	GetFeed(userIDs []int, offset, limit int, posts *[]Post) error
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

func (r *repository) GetFeed(userIDs []int, offset, limit int, posts *[]Post) error {
	return r.db.
		Where("user_id IN ?", userIDs).
		Order("created_at desc").
		Offset(offset).
		Limit(limit).
		Find(posts).Error
}
