package serve

import (
	"errors"
	"fmt"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli/serve/controller/dto"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/database"
	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"gorm.io/gorm"
)

func GetTrainerProfile(id uint) (models.Trainer, error) {
	var trainer models.Trainer
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("id = ?", id).First(&trainer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainer{}, nil
		}
		return models.Trainer{}, err
	}
	return trainer, nil
}

func GetTrainerById(id uint) (models.Trainer, error) {
	var trainer models.Trainer
	if err := database.DB.Preload("User").Preload("ActiveDays").Where("id = ?", id).First(&trainer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Trainer{}, nil
		}
		return models.Trainer{}, err
	}
	return trainer, nil
}

func EditTrainerProfile(id uint64, trainer dto.TrainerEdit) (models.Trainer, error) {
	user, _ := GetUserById(id)

	activeDays := models.ActiveDays{
		Monday:    trainer.ActiveDays[0],
		Tuesday:   trainer.ActiveDays[1],
		Wednesday: trainer.ActiveDays[2],
		Thursday:  trainer.ActiveDays[3],
		Friday:    trainer.ActiveDays[4],
		Saturday:  trainer.ActiveDays[5],
		Sunday:    trainer.ActiveDays[6],
	}
	if err := database.DB.Create(&activeDays).Error; err != nil {
		return models.Trainer{}, err
	}

	newTrainer := models.Trainer{
		UserID:          user.ID,
		UserName:        trainer.UserName,
		Status:          trainer.Status,
		CoachExperience: trainer.CoachExperience,
		Contact:         trainer.Contact,
		Language:        trainer.Language,
		Country:         trainer.Country,
		Sport:           trainer.Sport,
		Achievements:    trainer.Achievements,
		Education:       trainer.Education,
		ActiveDaysID:    activeDays.ID,
	}

	user.InfoID = fmt.Sprintf("%d", newTrainer.ID)

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingTrainer models.Trainer
	result := tx.Where(models.Trainee{UserID: newTrainer.UserID}).First(&existingTrainer)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := tx.Create(&newTrainer).Error; err != nil {
			tx.Rollback()
			return models.Trainer{}, err
		}
	} else if result.Error == nil {
		if err := tx.Model(&existingTrainer).Updates(&newTrainer).Error; err != nil {
			tx.Rollback()
			return models.Trainer{}, err
		}
		newTrainer = existingTrainer
	} else {
		tx.Rollback()
		return models.Trainer{}, result.Error
	}

	user.InfoID = fmt.Sprintf("%d", newTrainer.ID)
	if err := tx.Model(&user).Updates(models.User{InfoID: user.InfoID}).Error; err != nil {
		tx.Rollback()
		return models.Trainer{}, err
	}

	tx.Commit()

	createdTrainer, err := GetTrainerById(newTrainer.ID)
	if err != nil {
		return models.Trainer{}, err
	}
	return createdTrainer, nil
}

func GetRequests(trainer models.Trainer) ([]models.Request, error) {
	var requests []models.Request
	for _, i := range trainer.RequestIDs {
		r, err := GetRequest(uint(i))
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}
	return requests, nil
}

func SetPrice(setPrice dto.TrainerSetPrice) (models.Request, error) {
	r, err := GetRequest(setPrice.RequestId)
	if err != nil {
		return models.Request{}, err
	}
	tx := database.DB.Begin()
	if setPrice.Rejected {
		r.Status = "TrainerRejected"
	} else {
		r.Price = setPrice.Price
		r.Status = "TraineePending"
	}
	if err := tx.Save(&r).Error; err != nil {
		tx.Rollback()
		return models.Request{}, err
	}
	tx.Commit()
	return r, nil
}
