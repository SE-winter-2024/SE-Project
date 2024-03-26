package dto

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint   `json:"id"`
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Age         uint   `json:"age" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" validate:"required,min=11"`
	InfoID      string `json:"info_id"`
	InfoType    string `json:"info_type"`
	Block       bool   `json:"block"`
	Wallet      uint64 `json:"wallet"`
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPasswordHash(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	return err
}

type SetProgram struct {
	Activity string `json:"activity,omitempty"`
	Count    uint   `json:"count,omitempty"`
	Weight   uint   `json:"weight,omitempty"`
	Time     uint   `json:"time,omitempty"`
}

type TrainerEdit struct {
	User            models.User `json:"user"`
	UserName        string      `json:"user_name" validate:"required"`
	Status          string      `json:"status" validate:"required"`
	CoachExperience uint        `json:"coach_experience" validate:"required"`
	Contact         string      `json:"contact" validate:"required"`
	Language        string      `json:"language" validate:"required"`
	Country         string      `json:"country" validate:"required"`
	Sport           string      `json:"sport" validate:"required"`
	Achievements    string      `json:"achievements" validate:"required"`
	Education       string      `json:"education" validate:"required"`
	ActiveDays      []bool      `json:"active_days" validate:"required"`
}

type TraineeEdit struct {
	User           models.User `json:"user"`
	Height         uint        `json:"height" validate:"required"`
	Weight         uint        `json:"weight" validate:"required"`
	Sports         string      `json:"sports" validate:"required"`
	UserName       string      `json:"user_name" validate:"required"`
	Status         string      `json:"status" validate:"required"`
	Contact        string      `json:"contact" validate:"required"`
	Language       string      `json:"language" validate:"required"`
	Country        string      `json:"country" validate:"required"`
	MedicalHistory string      `json:"medicalHistory" validate:"required"`
	ActiveDays     []bool      `json:"active_days" validate:"required"`
}