package article

import "time"

type Status string

const (
	StatusPublish Status = "publish"
	StatusDraft   Status = "draft"
	StatusThrash  Status = "thrash"
)

type Article struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Content     string    `db:"content" json:"content"`
	Category    string    `db:"category" json:"category"`
	CreatedDate time.Time `db:"created_date" json:"created_date"`
	UpdatedDate time.Time `db:"updated_date" json:"updated_date"`
	Status      Status    `db:"status" json:"status"`
}
