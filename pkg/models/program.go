package models

import (
	"time"

	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	TraineeName  string
	TraineeID    uint
	TrainerID    uint
	Date         time.Time
	Status       string
	Description  string
	Price        uint
	ActiveDaysID uint
	ActiveDays   ActiveDays
}

type ActiveDays struct {
	gorm.Model
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	Sunday    bool
}

type TrainingProgram struct {
	gorm.Model
	TraineeID    uint
	Title        string
	Description  string
	StartDate    time.Time
	EndDate      time.Time
	ActivityDays []ActiveDays
}

type SportActivity struct {
	gorm.Model
	OrderNumber   uint
	ExpectedValue uint
	Value         uint
	Status        string
	Sport
}

type Sport struct {
	gorm.Model
	Title       string
	Description string
	Video       Media
}

type Media struct {
	gorm.Model
	Path      string
	Name      string
	Type      string
	Size      uint
	MediaId   string
	MediaType string
}
