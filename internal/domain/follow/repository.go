package follow

import "gorm.io/gorm"

type Repository interface {
	Follow(f *Follow) error
	Unfollow(followerID, followeeID string) error
	GetFollowingIDs(followerID string) ([]string, error)
	IsFollowing(followerID, followeeID string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Follow(f *Follow) error {
	return r.db.Create(f).Error
}

func (r *repository) Unfollow(followerID, followeeID string) error {
	return r.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&Follow{}).Error
}

func (r *repository) GetFollowingIDs(followerID string) ([]string, error) {
	var ids []string
	err := r.db.Model(&Follow{}).
		Where("follower_id = ?", followerID).
		Pluck("followee_id", &ids).Error
	return ids, err
}

func (r *repository) IsFollowing(followerID, followeeID string) (bool, error) {
	var count int64
	err := r.db.Model(&Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error
	return count > 0, err
}
