package post

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Content   string    `gorm:"type:varchar(200);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
