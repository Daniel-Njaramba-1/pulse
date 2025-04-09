package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
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

func (admin *Admin) FeedGetId() *int {
	return &admin.Id
}

func (admin *Admin) FeedCreateQuery() string {
	return `
		INSERT INTO admins (username, email, password_hash, is_active)
		VALUES (:username, :email, :password_hash, :is_active)
		RETURNING id
	`
}

func (admin *Admin) FeedGetByIdQuery() string {
	return `
		SELECT id, username, email, is_active, created_at, updated_at
		FROM admins
		WHERE id = $1
	`
}

func (admin *Admin) FeedGetAllQuery() string {
	return `
		SELECT id, username, email, is_active, created_at, updated_at
		FROM admins
		ORDER BY id ASC
	`
}

func (admin *Admin) FeedUpdateDetailsQuery() string {
	return `
		UPDATE admins
		SET username = :username,
			email = :email
		WHERE id = :id
	`
}

func (admin *Admin) FeedDeactivateQuery() string {
	return `
		UPDATE admins
		SET is_active = FALSE
		WHERE id = :id
	`
}

func (admin *Admin) FeedReactivateQuery() string {
	return `
		UPDATE admins
		SET is_active = TRUE
		WHERE id = :id
	`
}

func (admin *Admin) FeedDeleteQuery() string {
	return `
		DELETE FROM admins
		WHERE id = :id
	`
}

func GetAdminByNameQuery(ctx context.Context, db *sqlx.DB, name string, admin *Admin) error {
	query := 
	`	SELECT id, username, email, is_active, created_at, updated_at
		FROM admins 
		WHERE username = $1
	`
	return db.Get(admin, query, name)
}

func SearchAdminByNameQuery(ctx context.Context, db *sqlx.DB, name string, admin *Admin) error {
	query := 
	`	SELECT id, username, email, is_active, created_at, updated_at
		FROM admins 
		WHERE username ILIKE $1
	`
	return db.Get(admin, query, name)
}

// for comparing with password attempts on login
func GetAdminPasswordQuery(ctx context.Context, db *sqlx.DB, username string, admin *Admin) error {
	query := 
	`	SELECT password_hash
		FROM admins 
		WHERE username = $1
	`
	return db.Get(admin, query, username)
}

// update password using ID
func UpdateAdminPasswordQuery(ctx context.Context, db *sqlx.DB, id int, admin *Admin) error {
	query := 
	`	UPDATE admins
		SET password_hash = :password_hash
		WHERE id = :id
	`
	_, err := db.NamedExec(query, admin)
	return err
}