package repo

import (
	"time"
)

type Category struct {
	Id			int			`db:"id" json:"id"`
	Name		string		`db:"name" json:"name"`
	Description	string		`db:"description" json:"description"`
	IsActive	bool		`db:"is_active" json:"is_active"`
	CreatedAt	time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt	time.Time	`db:"updated_at" json:"updated_at"`
}
