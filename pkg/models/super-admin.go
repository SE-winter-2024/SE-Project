package models

import (
	"gorm.io/gorm"
)

type SuperAdmin struct {
	gorm.Model
	User   User
	UserID uint
}
