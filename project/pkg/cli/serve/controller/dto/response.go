package dto

type UserResponse struct {
	Email       string `json:"email,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Password    string `json:"password,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Type        string `json:"type,omitempty"`
	InfoID      string `json:"info_id,omitempty"`
	InfoType    string `json:"info_type,omitempty"`
	Block       bool   `json:"block,omitempty"`
	Wallet      uint64 `json:"wallet,omitempty"`
}

type TrainerResponse struct {
	TrainerProfileCard `json:"trainer_profile_card"`
	Sports             string `json:"sports,omitempty"`
	Achievements       string `json:"achievements,omitempty"`
	Education          string `json:"education,omitempty"`
}

type TrainerProfileCard struct {
	UserName        string `json:"user_name,omitempty"`
	Email           string `json:"email,omitempty"`
	Status          string `json:"status,omitempty"`
	Role            string `json:"role,omitempty"`
	CoachExperience uint   `json:"coach_experience,omitempty"`
	Contact         string `json:"contact,omitempty"`
	Language        string `json:"language,omitempty"`
	Country         string `json:"country,omitempty"`
}
