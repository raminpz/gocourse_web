package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	FirstName string `json:"first_name" gorm:"type:char(50);not null"`
	LastName  string `json:"last_name" gorm:"type:char(50);not null"`
	Email     string `json:"email" gorm:"type:char(50);not null; unique"`
	Phone     string `json:"phone" gorm:"type:char(11);not null; unique"`
	Password  string `json:"password" gorm:"type:char(50);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
