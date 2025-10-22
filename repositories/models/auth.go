package models

import (
	"time"

	"github.com/google/uuid"
)

type BannedLists struct {
	BannedListId uuid.UUID `gorm:"primaryKey;autoIncrement" json:"banned_list_id"`
	Token        string    `gorm:"type:varchar(300);unique" json:"token"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type LoggedLists struct {
	LoggedListId uuid.UUID `gorm:"primaryKey;autoIncrement" json:"logged_list_id"`
	Token        string    `gorm:"type:varchar(300);unique" json:"token"`
	UserEmail    string    `gorm:"type:varchar(300)" json:"user_email"`
	UserId       uuid.UUID `json:"user_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
