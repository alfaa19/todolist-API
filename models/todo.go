package models

import "time"

type Todo struct {
	TodoID          uint      `gorm:"primaryKey" json:"id"`
	ActivityGroupID uint      `json:"activity_group_id" validate:"required,number,gt=0"`
	Title           string    `gorm:"type:varchar(255)" json:"title" validate:"required"`
	IsActive        bool      `gorm:"type:bool" json:"is_active" validate:"required,boolean"`
	Priority        string    `gorm:"type:varchar(55)" json:"priority" validate:"required"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
