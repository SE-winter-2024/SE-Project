package models

import (
	"gorm.io/gorm"
)

type Trainee struct {
	gorm.Model
	User           User
	UserID         uint
	Height         uint
	Weight         uint
	Sports         string
	UserName       string
	Status         string
	Contact        string
	Language       string
	Country        string
	MedicalHistory string
	ActiveDays     ActiveDays
	ActiveDaysID   uint
	TrainerID      uint
	Trainer        Trainer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RequestID      uint
}
