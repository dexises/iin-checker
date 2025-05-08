package service

import (
	"context"
	"errors"
	databases "github.com/dexises/iin-checker/internal/databases/drivers/mongo"
	"github.com/dexises/iin-checker/internal/models"
	"time"
)

type PersonService interface {
	// ValidateIIN accepts a string IIN and returns date of birth, gender, ok flag, and error.
	ValidateIIN(iin string) (time.Time, string, bool, error)

	Create(ctx context.Context, req models.CreatePersonRequest) (string, error)
	Get(ctx context.Context, iin string) (models.PersonDTO, error)
	FindByName(ctx context.Context, namePart string) ([]models.PersonDTO, error)
}

// personService is a concrete implementation of PersonService.
type personService struct {
	repo databases.PersonRepository
}

// NewPersonService constructs a PersonService with given repository.
func NewPersonService(repo databases.PersonRepository) PersonService {
	return &personService{repo: repo}
}

func (s *personService) FindByName(ctx context.Context, namePart string) ([]models.PersonDTO, error) {
	persons, err := s.repo.FindByName(ctx, namePart)
	if err != nil {
		return nil, err
	}
	out := make([]models.PersonDTO, len(persons))
	for i, p := range persons {
		out[i] = models.PersonDTO{IIN: p.IIN, Name: p.Name, Phone: p.Phone}
	}
	return out, nil
}

// ValidateIIN implements full IIN validation: format, date, gender, checksum.
func (s *personService) ValidateIIN(iin string) (date time.Time, gender string, ok bool, err error) {
	// Length check
	if len(iin) != 12 {
		err = errors.New("IIN должен содержать ровно 12 цифр")
		return
	}

	// Parse digits
	digits := make([]int, 12)
	for i, r := range iin {
		if r < '0' || r > '9' {
			err = errors.New("IIN должен состоять только из цифр")
			return
		}
		digits[i] = int(r - '0')
	}

	// Extract date components YYMMDD
	yy := digits[0]*10 + digits[1]
	mm := digits[2]*10 + digits[3]
	dd := digits[4]*10 + digits[5]

	// Determine century and gender from 7th digit
	switch digits[6] {
	case 1:
		date = time.Date(1800+yy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
		gender = "male"
	case 2:
		date = time.Date(1800+yy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
		gender = "female"
	case 3:
		date = time.Date(1900+yy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
		gender = "male"
	case 4:
		date = time.Date(1900+yy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
		gender = "female"
	case 5:
		date = time.Date(2000+yy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
		gender = "male"
	case 6:
		date = time.Date(2000+yy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
		gender = "female"
	default:
		err = errors.New("некорректная 7-я цифра: невозможно определить век/пол")
		return
	}

	// Validate actual date
	if date.Month() != time.Month(mm) || date.Day() != dd {
		err = errors.New("некорректная дата рождения в ИИН")
		return
	}

	// Checksum calculation (first weights)
	w1 := [11]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	sum := 0
	for i := 0; i < 11; i++ {
		sum += digits[i] * w1[i]
	}
	c := sum % 11
	if c == 10 {
		// Second weights if first mod == 10
		w2 := [11]int{3, 4, 5, 6, 7, 8, 9, 10, 11, 1, 2}
		sum = 0
		for i := 0; i < 11; i++ {
			sum += digits[i] * w2[i]
		}
		c = sum % 11
		if c == 10 {
			ok = false
			return
		}
	}
	// Final check digit
	ok = (c == digits[11])
	return
}

// Create validates and persists a person record.
func (s *personService) Create(ctx context.Context, req models.CreatePersonRequest) (string, error) {
	// Validate IIN
	if _, _, ok, err := s.ValidateIIN(req.IIN); err != nil || !ok {
		if err != nil {
			return "", err
		}
		return "", errors.New("invalid IIN")
	}
	// Check uniqueness
	exists, err := s.repo.Exists(ctx, req.IIN)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("person with this IIN already exists")
	}
	// Persist entity
	en := databases.Person{IIN: req.IIN, Name: req.Name, Phone: req.Phone}
	if err := s.repo.Create(ctx, en); err != nil {
		return "", err
	}
	return req.IIN, nil
}

// Get retrieves a person by IIN.
func (s *personService) Get(ctx context.Context, iin string) (models.PersonDTO, error) {
	// Validate IIN
	if _, _, ok, err := s.ValidateIIN(iin); err != nil || !ok {
		if err != nil {
			return models.PersonDTO{}, err
		}
		return models.PersonDTO{}, errors.New("invalid IIN")
	}
	p, err := s.repo.Get(ctx, iin)
	if err != nil {
		return models.PersonDTO{}, err
	}
	return models.PersonDTO{IIN: p.IIN, Name: p.Name, Phone: p.Phone}, nil
}
