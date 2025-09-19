package post

import "time"

type Post struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `gorm:"not null;index"`
	Content   string    `gorm:"type:varchar(200);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
