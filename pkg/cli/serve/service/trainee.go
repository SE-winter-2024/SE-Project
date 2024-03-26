package serve

import (
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func GetTraineeProfile(id uint) (models.Trainee, error) {
	var trainee models.Trainee
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("id = ?", id).First(&trainee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainee{}, nil
		}
		return models.Trainee{}, err
	}
	return trainee, nil
}

func GetTraineeById(id uint) (models.Trainee, error) {
	var trainee models.Trainee
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("id = ?", id).First(&trainee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainee{}, nil
		}
		return models.Trainee{}, err
	}
	return trainee, nil
}

func EditTraineeProfile(id uint64, trainee dto.TraineeEdit) (models.Trainee, error) {
	user, err := GetUserById(id)

	activeDays := models.ActiveDays{
		Monday:    trainee.ActiveDays[0],
		Tuesday:   trainee.ActiveDays[1],
		Wednesday: trainee.ActiveDays[2],
		Thursday:  trainee.ActiveDays[3],
		Friday:    trainee.ActiveDays[4],
		Saturday:  trainee.ActiveDays[5],
		Sunday:    trainee.ActiveDays[6],
	}

	if err := database.DB.Create(&activeDays).Error; err != nil {
		return models.Trainee{}, err
	}

	newTrainee := models.Trainee{
		UserID:         user.ID,
		Height:         trainee.Height,
		Weight:         trainee.Weight,
		UserName:       trainee.UserName,
		Status:         trainee.Status,
		Contact:        trainee.Contact,
		Language:       trainee.Language,
		Country:        trainee.Country,
		Sports:         trainee.Sports,
		MedicalHistory: trainee.MedicalHistory,
		ActiveDaysID:   activeDays.ID,
	}

	user.InfoID = fmt.Sprintf("%d", newTrainee.ID)

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingTrainee models.Trainee
	result := tx.Where(models.Trainee{UserID: newTrainee.UserID}).First(&existingTrainee)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := tx.Create(&newTrainee).Error; err != nil {
			tx.Rollback()
			return models.Trainee{}, err
		}
	} else if result.Error == nil {
		if err := tx.Model(&existingTrainee).Updates(&newTrainee).Error; err != nil {
			tx.Rollback()
			return models.Trainee{}, err
		}
		newTrainee = existingTrainee
	} else {
		tx.Rollback()
		return models.Trainee{}, result.Error
	}

	user.InfoID = fmt.Sprintf("%d", newTrainee.ID)
	if err := tx.Model(&user).Updates(models.User{InfoID: user.InfoID}).Error; err != nil {
		tx.Rollback()
		return models.Trainee{}, err
	}

	tx.Commit()

	createdTrainee, err := GetTraineeById(newTrainee.ID)
	if err != nil {
		return models.Trainee{}, err
	}
	return createdTrainee, nil
}
