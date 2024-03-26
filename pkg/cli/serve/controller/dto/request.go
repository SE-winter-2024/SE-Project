package dto

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID          uint   `json:"id"`
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" validate:"required,min=11"`
	Type        string `json:"type" validate:"required"`
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
