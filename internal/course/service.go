package course

import (
	"github.com/raminpz/gocourse_web/internal/domain"
	"log"
	"time"
)

type (
	Filters struct {
		Name string
	}
	Service interface {
		Create(name, startDate, endDate string) (*domain.Course, error)
		GetAll(filters Filters, limit, offset int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Update(id string, name, startDate, endDate *string) error
		Delete(id string) error
		Count(filters Filters) (int, error)
	}
	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(repo Repository, logger *log.Logger) Service {
	return &service{
		log:  logger,
		repo: repo,
	}
}

func (s service) Create(name, startDate, endDate string) (*domain.Course, error) {

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println("Error parsing start date:", err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println("Error parsing end date:", err)
		return nil, err
	}
	course := &domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}
	if err := s.repo.Create(course); err != nil {
		return nil, err
	}
	return course, nil

}

func (s service) GetAll(filters Filters, limit, offset int) ([]domain.Course, error) {
	courses, err := s.repo.GetAll(filters, limit, offset)
	if err != nil {
		s.log.Println("Error getting courses:", err)
		return nil, err
	}
	return courses, nil
}

func (s service) Get(id string) (*domain.Course, error) {
	course, err := s.repo.Get(id)
	if err != nil {
		s.log.Println("Error getting course:", err)
		return nil, err
	}
	return course, nil
}

func (s service) Update(id string, name, startDate, endDate *string) error {
	var startDateParsed, endDateParsed *time.Time
	if startDate != nil {
		parsed, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Println("Error parsing start date:", err)
			return err
		}
		startDateParsed = &parsed
	}
	if endDate != nil {
		parsed, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Println("Error parsing end date:", err)
			return err
		}
		endDateParsed = &parsed
	}
	return s.repo.Update(id, name, startDateParsed, endDateParsed)
}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Count(filters Filters) (int, error) {
	count, err := s.repo.Count(filters)
	if err != nil {
		s.log.Println("Error counting courses:", err)
		return 0, err
	}
	return count, nil
}
