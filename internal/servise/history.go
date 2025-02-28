package servise

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"systemgroup.net/bootcamp/go/v1/shell/internal/models"
	database "systemgroup.net/bootcamp/go/v1/shell/internal/storage"
)

func AddCommandHistory(username, command string) error {
    db := database.GetDB()

    var history models.CommandHistory
    result := db.Where("username = ? AND command = ?", username, command).First(&history)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        history = models.CommandHistory{
            Username:  username,
            Command:   command,
            CreatedAt: time.Now(),
            Count:     1,
        }
        return db.Create(&history).Error
    } else if result.Error != nil {
        return result.Error
    }

    history.Count += 1
    history.CreatedAt = time.Now()
    return db.Save(&history).Error
}

func GetCommandHistory(username string) ([]models.CommandHistory, error) {
    db := database.GetDB()

    var history []models.CommandHistory
    result := db.
        Where("username = ?", username).
        Select("command, count").
        Order("count DESC, created_at DESC").
        Find(&history)

    if result.Error != nil {
        return nil, result.Error
    }

    return history, nil
}