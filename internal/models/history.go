package models

import "time"

type CommandHistory struct {
    Username  string    `gorm:"primaryKey;column:username"`
    Command   string    `gorm:"primaryKey;column:command"`
    CreatedAt time.Time `gorm:"column:created_at"`
    Count     int       `gorm:"column:count"`
}

