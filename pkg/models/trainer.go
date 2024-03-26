package models

import "gorm.io/gorm"

type Trainer struct {
	gorm.Model
	User            User
	UserID          uint
	UserName        string
	Status          string
	CoachExperience uint
	Contact         string
	Language        string
	Country         string
	Sport           string
	Achievements    string
	Education       string
	ActiveDays      ActiveDays
	ActiveDaysID    uint
	Trainees        []Trainee
	Requests        []Request
}
