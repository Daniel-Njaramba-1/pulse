package adminSvc

import (
	"context"
	"errors"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/generics"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/hashing"
	"github.com/jmoiron/sqlx"
)

type Authentication struct {
	db *sqlx.DB
}

func NewAuthentication(db *sqlx.DB) *Authentication {
	return &Authentication{db: db}
}

func (a *Authentication) RegisterAdmin(ctx context.Context, admin *repo.Admin) (string, *repo.Admin, error) {
	if admin.Username == "" || admin.Email == "" || admin.Password == "" {
		return "", nil, errors.New("missing required fields")
	}

	hashedPassword, err := hashing.HashPassword(admin.Password)
	if err != nil {
		return "", nil, err
	}
	admin.PasswordHash = hashedPassword
	admin.Password = ""

	_, err = generics.CreateModel(ctx, a.db, admin)
	if err != nil {
		return "", nil, err
	}

	token, err := CreateAdminToken(admin.Username)
	if err != nil {
		return "", nil, err
	}

	return token, admin, nil
}

func (a *Authentication) LoginAdmin(ctx context.Context, username string, password string) (string, *repo.Admin, error) {
	var admin repo.Admin
	err := repo.GetAdminByNameQuery(ctx, a.db, username, &admin)
	if err != nil {
		return "", nil, err
	}
	err = repo.GetAdminPasswordQuery(ctx, a.db, username, &admin)
	if err != nil {
		return "", nil, err
	}

	if !hashing.VerifyPassword(password, admin.PasswordHash) {
		return "", nil, errors.New("invalid password")
	}

	token, err := CreateAdminToken(admin.Username)
	if err != nil {
		return "", nil, err
	}

	return token, &admin, nil
}

func (a *Authentication) ResetAdminPassword(ctx context.Context) {

}

func (a *Authentication) ResetCustomerPassword(ctx context.Context) {

}
