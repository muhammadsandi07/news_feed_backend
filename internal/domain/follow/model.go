package follow

import "time"

type Follow struct {
	FollowerID int       `gorm:"primaryKey;index" json:"follower_id"`
	FolloweeID int       `gorm:"primaryKey;index" json:"followee_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
