package article

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct{ Fields map[string]string }

func (e *ValidationError) Error() string { return "request validation failed" }

type Service struct {
	repository Repository
	validator  *validator.Validate
}

func NewService(repository Repository) *Service {
	v := validator.New()
	_ = v.RegisterValidation("article_status", func(field validator.FieldLevel) bool {
		return ValidStatus(Status(field.Field().String()))
	})
	return &Service{repository: repository, validator: v}
}

func ValidStatus(status Status) bool {
	return status == StatusPublish || status == StatusDraft || status == StatusThrash
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Article, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	input.Category = strings.TrimSpace(input.Category)
	if err := s.validate(input); err != nil {
		return Article{}, err
	}
	return s.repository.Create(ctx, input)
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]Article, error) {
	if filter.Limit < 1 || filter.Limit > 100 || filter.Offset < 0 {
		return nil, &ValidationError{Fields: map[string]string{"pagination": "limit must be between 1 and 100 and offset must be zero or greater"}}
	}
	if filter.Status != nil && !ValidStatus(*filter.Status) {
		return nil, &ValidationError{Fields: map[string]string{"status": "status must be one of publish, draft, or thrash"}}
	}
	return s.repository.List(ctx, filter)
}

func (s *Service) GetByID(ctx context.Context, id int64) (Article, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (Article, error) {
	current, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return Article{}, err
	}
	if input.Title != nil {
		current.Title = strings.TrimSpace(*input.Title)
	}
	if input.Content != nil {
		current.Content = strings.TrimSpace(*input.Content)
	}
	if input.Category != nil {
		current.Category = strings.TrimSpace(*input.Category)
	}
	if input.Status != nil {
		current.Status = *input.Status
	}

	validated := CreateInput{Title: current.Title, Content: current.Content, Category: current.Category, Status: current.Status}
	if err := s.validate(validated); err != nil {
		return Article{}, err
	}
	return s.repository.Update(ctx, current)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) validate(input CreateInput) error {
	if err := s.validator.Struct(input); err != nil {
		fields := make(map[string]string)
		var validationErrors validator.ValidationErrors
		if ok := errors.As(err, &validationErrors); !ok {
			return fmt.Errorf("validate article: %w", err)
		}
		for _, fieldError := range validationErrors {
			name := strings.ToLower(fieldError.Field())
			switch fieldError.Tag() {
			case "required":
				fields[name] = name + " is required"
			case "min":
				fields[name] = fmt.Sprintf("%s must contain at least %s characters", name, fieldError.Param())
			case "max":
				fields[name] = fmt.Sprintf("%s must contain at most %s characters", name, fieldError.Param())
			case "article_status":
				fields[name] = "status must be one of publish, draft, or thrash"
			}
		}
		return &ValidationError{Fields: fields}
	}
	return nil
}
