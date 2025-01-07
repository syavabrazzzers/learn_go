package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);unique_index" json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" default:"false"`
}

type UserRecoveryCodes struct {
	UserId int    `gorm:"index;primarykey"`
	User   User   `gorm:"constraint:OnDelete:CASCADE"`
	Codes  []byte `gorm:"type:jsonb" json:"codes"`
}
