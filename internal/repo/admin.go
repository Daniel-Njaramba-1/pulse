package repo

import (
	"time"
)

type Admin struct {
    Id				int			`db:"id" json:"id"`
    Username		string		`db:"username" json:"username"`
    Email			string		`db:"email" json:"email"`
    PasswordHash	string		`db:"password_hash" json:"-"`
    Password        string      `db:"-" json:"password"` 
    IsActive		bool		`db:"is_active" json:"is_active"`
    CreatedAt		time.Time	`db:"created_at" json:"created_at"`
    UpdatedAt		time.Time	`db:"updated_at" json:"updated_at"`
}
