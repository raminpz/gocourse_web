package enrollment

import (
	"github.com/raminpz/gocourse_web/internal/domain"
	"gorm.io/gorm"
	"log"
)

type (
	Repository interface {
		Create(enroll *domain.Enrollment) error
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(db *gorm.DB, logger *log.Logger) Repository {
	return &repo{
		db:  db,
		log: logger,
	}
}

func (r *repo) Create(enroll *domain.Enrollment) error {
	if err := r.db.Create(enroll).Error; err != nil {
		r.log.Println("error: %v", err)
		return err
	}
	r.log.Println("enrollment created with id:", enroll.ID)
	return nil
}
