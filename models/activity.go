package models

import (
	"time"
)

type Activity struct {
	ActivityID uint      `gorm:"primaryKey" json:"id"`
	Title      string    `gorm:"not null;type:varchar(255)" json:"title" validate:"required"`
	Email      string    `gorm:"type:varchar(255)" json:"email" validate:"required,email"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
