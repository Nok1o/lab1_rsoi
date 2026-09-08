package person

import (
	"context"
	"strings"

	"github.com/Nok1o/lab1_rsoi/backend/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, person domain.Person) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.Person, error)
	List(ctx context.Context) ([]domain.Person, error)
	Update(ctx context.Context, person domain.Person) error
	Delete(ctx context.Context, id int64) error
}

type CreateInput struct {
	Name    string
	Age     int32
	Address string
	Work    string
}

type UpdateInput struct {
	Name    *string
	Age     *int32
	Address *string
	Work    *string
}

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (int64, error) {
	if err := validateCreate(input); err != nil {
		return 0, err
	}

	return s.repository.Create(ctx, domain.Person{
		Name:    strings.TrimSpace(input.Name),
		Age:     input.Age,
		Address: input.Address,
		Work:    input.Work,
	})
}

func (s *Service) GetByID(ctx context.Context, id int64) (domain.Person, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]domain.Person, error) {
	return s.repository.List(ctx)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (domain.Person, error) {
	if err := validateUpdate(input); err != nil {
		return domain.Person{}, err
	}

	current, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Person{}, err
	}
	if input.Name != nil {
		current.Name = strings.TrimSpace(*input.Name)
	}
	if input.Age != nil {
		current.Age = *input.Age
	}
	if input.Address != nil {
		current.Address = *input.Address
	}
	if input.Work != nil {
		current.Work = *input.Work
	}

	if err := s.repository.Update(ctx, current); err != nil {
		return domain.Person{}, err
	}
	return current, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func validateCreate(input CreateInput) error {
	errors := make(map[string]string)
	if strings.TrimSpace(input.Name) == "" {
		errors["name"] = "must not be empty"
	}
	if input.Age < 0 {
		errors["age"] = "must be greater than or equal to zero"
	}
	if len(errors) != 0 {
		return &ValidationError{Fields: errors}
	}
	return nil
}

func validateUpdate(input UpdateInput) error {
	errors := make(map[string]string)
	if input.Name == nil && input.Age == nil && input.Address == nil && input.Work == nil {
		errors["body"] = "must contain at least one field"
	}
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		errors["name"] = "must not be empty"
	}
	if input.Age != nil && *input.Age < 0 {
		errors["age"] = "must be greater than or equal to zero"
	}
	if len(errors) != 0 {
		return &ValidationError{Fields: errors}
	}
	return nil
}
