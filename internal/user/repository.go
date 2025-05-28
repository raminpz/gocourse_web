package user

import (
	"fmt"
	"github.com/raminpz/gocourse_web/internal/domain"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *domain.User) error
	GetAll(filters Filters, limit, offset int) ([]domain.User, error)
	GetByID(id string) (*domain.User, error)
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(filters Filters) (int, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (r *repo) Create(user *domain.User) error {

	if err := r.db.Create(user).Error; err != nil {
		r.log.Println(err)
		return err
	}
	r.log.Println("User created with id: ", user.ID)
	return nil
}

func (r *repo) GetAll(filters Filters, limit, offset int) ([]domain.User, error) {
	var user []domain.User
	tx := r.db.Model(&user)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	tx = applyFilters(tx, filters)
	result := tx.Order("created_at desc").Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil

}

func (r *repo) GetByID(id string) (*domain.User, error) {
	user := domain.User{ID: id}
	result := r.db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *repo) Delete(id string) error {
	user := domain.User{ID: id}
	result := r.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repo) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {
	values := make(map[string]interface{})
	if firstName != nil {
		values["first_name"] = firstName
	}
	if lastName != nil {
		values["last_name"] = lastName
	}
	if email != nil {
		values["email"] = email
	}
	if phone != nil {
		values["phone"] = phone
	}
	result := r.db.Model(&domain.User{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		// Usamos LOWER tanto en la columna como en el valor para asegurar una búsqueda case-insensitive
		tx = tx.Where("LOWER(first_name) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", filters.FirstName))
	}
	if filters.LastName != "" {
		tx = tx.Where("LOWER(last_name) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", filters.LastName))
	}
	return tx
}

func (r *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := r.db.Model(&domain.User{})
	tx = applyFilters(tx, filters)
	result := tx.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}
