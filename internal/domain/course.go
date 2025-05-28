package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Course struct {
	ID        string    `json:"id" gorm:"type:char(36);not null;primaryKey;unique"`
	Name      string    `json:"name" gorm:"type:char(50);not null"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	User      *User     `gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
