package db

import (
	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	ImageData []byte `json:"image_data"` // Store image data
	ImagePath string `json:"image_path"` // Optional: Store image path if needed
}
