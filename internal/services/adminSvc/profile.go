package adminSvc

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Profile struct {
	db *sqlx.DB
}

func NewProfile(db *sqlx.DB) *Profile {
	return &Profile{db: db}
}

func EditProfile(ctx context.Context, db *sqlx.DB) {

}

func DeactivateProfile(ctx context.Context, db *sqlx.DB) {

}
