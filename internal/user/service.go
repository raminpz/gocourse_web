package user

import "log"

type (
	Filters struct {
		FirstName string
		LastName  string
	}
	Service interface {
		Create(firstName, lastName, email, phone, password string) (*User, error)
		Get(id string) (*User, error)
		GetAll(filters Filters, limit, offset int) ([]User, error)
		Delete(id string) error
		Update(id string, firstName *string, lastName *string, email *string, phone *string, password *string) error
		Count(filters Filters) (int, error)
	}
	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(firstName, lastName, email, phone, password string) (*User, error) {
	s.log.Println("Create user service")
	user := User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Password:  password,
	}
	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s service) GetAll(filters Filters, limit, offset int) ([]User, error) {
	users, err := s.repo.GetAll(filters, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s service) Get(id string) (*User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Update(id string, firstName *string, lastName *string, email *string, phone *string, password *string) error {
	return s.repo.Update(id, firstName, lastName, email, phone, password)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
