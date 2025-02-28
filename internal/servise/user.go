package servise

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"systemgroup.net/bootcamp/go/v1/shell/internal/models"
	database "systemgroup.net/bootcamp/go/v1/shell/internal/storage"
)

var (
	ErrUserExists = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)


func CreateUser(user *models.User) error {
	db := database.GetDB()
	var existingUser models.User

	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func GetUser(inUser *models.User) (*models.User,error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("username = ?", inUser.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user.Password != inUser.Password {
		return nil, fmt.Errorf("wrong password")
	}


	return &user, nil
}