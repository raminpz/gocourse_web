package enrollment

import (
	"errors"
	"github.com/raminpz/gocourse_web/internal/course"
	"github.com/raminpz/gocourse_web/internal/domain"
	"github.com/raminpz/gocourse_web/internal/user"
	"log"
)

type (
	Service interface {
		Create(userID, courseID string) (*domain.Enrollment, error)
	}
	service struct {
		log       *log.Logger
		userSrv   user.Service
		courseSrv course.Service
		repo      Repository
	}
)

func NewService(repo Repository, logger *log.Logger, userSrv user.Service, courseSrv course.Service) Service {
	return &service{
		log:       logger,
		userSrv:   userSrv,
		courseSrv: courseSrv,
		repo:      repo,
	}
}

func (s service) Create(userID, courseID string) (*domain.Enrollment, error) {
	enroll := &domain.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   "P",
	}
	if _, err := s.userSrv.Get(enroll.UserID); err != nil {
		return nil, errors.New("user id does not exist: " + enroll.UserID)
	}
	if _, err := s.courseSrv.Get(enroll.CourseID); err != nil {
		return nil, errors.New("course id does not exist: " + enroll.CourseID)
	}

	if err := s.repo.Create(enroll); err != nil {
		s.log.Println("error creating enrollment:", err)
		return nil, err
	}

	s.log.Println("enrollment created with id:", enroll.ID)
	return enroll, nil
}
