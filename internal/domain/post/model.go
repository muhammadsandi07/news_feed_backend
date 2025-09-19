package post

import "time"

type Post struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"not null;index:idx_user_created" json:"userid"`
	Content   string    `gorm:"type:varchar(200);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
