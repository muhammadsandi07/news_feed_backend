package follow

import "time"

type Follow struct {
	FollowerID string    `gorm:"primaryKey"`
	FolloweeID string    `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
