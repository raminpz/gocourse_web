package course

import (
	"fmt"
	"github.com/raminpz/gocourse_web/internal/domain"
	"gorm.io/gorm"
	"log"
	"time"
)

type (
	Repository interface {
		Create(course *domain.Course) error
		GetAll(filters Filters, limit, offset int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Update(id string, name *string, startDate *time.Time, endDate *time.Time) error
		Delete(id string) error
		Count(filters Filters) (int, error)
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

func (r *repo) Create(course *domain.Course) error {
	if err := r.db.Create(course).Error; err != nil {
		r.log.Println("Error creating course:", err)
		return err
	}
	r.log.Println("Course created with id:", course.ID)
	return nil
}

func (r *repo) GetAll(filters Filters, limit, offset int) ([]domain.Course, error) {
	var courses []domain.Course
	tx := r.db.Model(&courses)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}
	return courses, nil
}

func (r *repo) Get(id string) (*domain.Course, error) {
	course := domain.Course{ID: id}
	result := r.db.First(&course)
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (r *repo) Delete(id string) error {
	course := domain.Course{ID: id}
	result := r.db.Delete(&course)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repo) Update(id string, name *string, startDate *time.Time, endDate *time.Time) error {
	values := make(map[string]interface{})
	if name != nil {
		values["name"] = *name
	}
	if startDate != nil {
		values["start_date"] = *startDate
	}
	if endDate != nil {
		values["end_date"] = *endDate
	}
	result := r.db.Model(&domain.Course{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := r.db.Model(&domain.Course{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil

}

func applyFilters(txt *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		txt = txt.Where("LOWER(name) LIKE (?)", fmt.Sprintf("%%%s%%", filters.Name))
	}
	return txt

}
