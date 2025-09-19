package follow

import "gorm.io/gorm"

type Repository interface {
	Follow(f *Follow) error
	Unfollow(followerID, followeeID int) error
	GetFollowingIDs(followerID int) ([]int, error)
	IsFollowing(followerID, followeeID int) (bool, error)
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

func (r *repository) Unfollow(followerID, followeeID int) error {
	return r.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&Follow{}).Error
}

func (r *repository) GetFollowingIDs(followerID int) ([]int, error) {
	var ids []int
	err := r.db.Model(&Follow{}).
		Where("follower_id = ?", followerID).
		Pluck("followee_id", &ids).Error
	return ids, err
}

func (r *repository) IsFollowing(followerID, followeeID int) (bool, error) {
	var count int64
	err := r.db.Model(&Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error
	return count > 0, err
}
