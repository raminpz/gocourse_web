package user

import (
	"github.com/google/uuid"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	GetAll() ([]User, error)
	GetByID(id string) (*User, error)
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone, password *string) error
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

func (r *repo) Create(user *User) error {
	user.ID = uuid.New().String()

	if err := r.db.Create(user).Error; err != nil {
		r.log.Println(err)
		return err
	}
	r.log.Println("User created with id: ", user.ID)
	return nil
}

func (r *repo) GetAll() ([]User, error) {
	var users []User
	result := r.db.Model(&users).Order("created_at desc").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil

}

func (r *repo) GetByID(id string) (*User, error) {
	user := User{ID: id}
	result := r.db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *repo) Delete(id string) error {
	user := User{ID: id}
	result := r.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repo) Update(id string, firstName *string, lastName *string, email *string, phone *string, password *string) error {
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
	if password != nil {
		values["password"] = password
	}
	result := r.db.Model(&User{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
