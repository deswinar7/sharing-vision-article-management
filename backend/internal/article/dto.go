package article

type CreateInput struct {
	Title    string `json:"title" validate:"required,min=20,max=200"`
	Content  string `json:"content" validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3,max=100"`
	Status   Status `json:"status" validate:"required,article_status"`
}

type UpdateInput struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	Category *string `json:"category"`
	Status   *Status `json:"status"`
}

type ListFilter struct {
	Limit  int
	Offset int
	Status *Status
}
