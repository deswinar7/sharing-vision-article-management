package article

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeRepository struct {
	articles map[int64]Article
	nextID   int64
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{articles: make(map[int64]Article), nextID: 1}
}

func (r *fakeRepository) Create(_ context.Context, input CreateInput) (Article, error) {
	now := time.Now().UTC()
	result := Article{ID: r.nextID, Title: input.Title, Content: input.Content, Category: input.Category, Status: input.Status, CreatedDate: now, UpdatedDate: now}
	r.articles[result.ID] = result
	r.nextID++
	return result, nil
}

func (r *fakeRepository) List(_ context.Context, filter ListFilter) ([]Article, error) {
	result := make([]Article, 0)
	for _, item := range r.articles {
		if filter.Status == nil || item.Status == *filter.Status {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *fakeRepository) GetByID(_ context.Context, id int64) (Article, error) {
	result, ok := r.articles[id]
	if !ok {
		return Article{}, ErrNotFound
	}
	return result, nil
}

func (r *fakeRepository) Update(_ context.Context, item Article) (Article, error) {
	if _, ok := r.articles[item.ID]; !ok {
		return Article{}, ErrNotFound
	}
	item.UpdatedDate = time.Now().UTC()
	r.articles[item.ID] = item
	return item, nil
}

func (r *fakeRepository) Delete(_ context.Context, id int64) error {
	if _, ok := r.articles[id]; !ok {
		return ErrNotFound
	}
	delete(r.articles, id)
	return nil
}

func validInput() CreateInput {
	return CreateInput{
		Title:    "A sufficiently descriptive article title",
		Content:  strings.Repeat("Useful article content for readers. ", 8),
		Category: "Engineering",
		Status:   StatusDraft,
	}
}

func TestCreateArticle(t *testing.T) {
	service := NewService(newFakeRepository())
	result, err := service.Create(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.ID != 1 || result.Status != StatusDraft {
		t.Fatalf("Create() result = %+v", result)
	}
}

func TestCreateValidation(t *testing.T) {
	tests := []struct {
		name  string
		alter func(*CreateInput)
		field string
	}{
		{name: "title too short", alter: func(input *CreateInput) { input.Title = "Too short" }, field: "title"},
		{name: "content too short", alter: func(input *CreateInput) { input.Content = "Too short" }, field: "content"},
		{name: "category too short", alter: func(input *CreateInput) { input.Category = "IT" }, field: "category"},
		{name: "invalid status", alter: func(input *CreateInput) { input.Status = "archived" }, field: "status"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := validInput()
			test.alter(&input)
			_, err := NewService(newFakeRepository()).Create(context.Background(), input)
			var validationError *ValidationError
			if !errors.As(err, &validationError) {
				t.Fatalf("Create() error = %v, want ValidationError", err)
			}
			if validationError.Fields[test.field] == "" {
				t.Fatalf("Create() fields = %v, want %q", validationError.Fields, test.field)
			}
		})
	}
}

func TestArticleNotFound(t *testing.T) {
	_, err := NewService(newFakeRepository()).GetByID(context.Background(), 99)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetByID() error = %v, want ErrNotFound", err)
	}
}

func TestPaginationValidation(t *testing.T) {
	tests := []ListFilter{{Limit: 0, Offset: 0}, {Limit: 101, Offset: 0}, {Limit: 10, Offset: -1}}
	for _, filter := range tests {
		_, err := NewService(newFakeRepository()).List(context.Background(), filter)
		var validationError *ValidationError
		if !errors.As(err, &validationError) {
			t.Fatalf("List(%+v) error = %v, want ValidationError", filter, err)
		}
	}
}

func TestUpdateArticleAndMoveToTrash(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository)
	created, err := service.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}
	newTitle := "An updated and still descriptive article title"
	updated, err := service.Update(context.Background(), created.ID, UpdateInput{Title: &newTitle})
	if err != nil || updated.Title != newTitle {
		t.Fatalf("Update() = %+v, %v", updated, err)
	}
	status := StatusThrash
	trashed, err := service.Update(context.Background(), created.ID, UpdateInput{Status: &status})
	if err != nil || trashed.Status != StatusThrash || trashed.Content != created.Content {
		t.Fatalf("trash update = %+v, %v", trashed, err)
	}
}
