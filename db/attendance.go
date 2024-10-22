package db

import (
	"time"

	"gorm.io/gorm"
)

// Attendance represents the attendance model
type Attendance struct {
	gorm.Model
	UserID      uint      `json:"user_id"`
	Timestamp   time.Time `json:"timestamp"`
	ImageData   []byte    `json:"image_data"`
	ImageFormat string    `json:"image_format"`
}
