package follow

import "time"

type Follow struct {
	FollowerID uint      `gorm:"primaryKey"`
	FolloweeID uint      `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
