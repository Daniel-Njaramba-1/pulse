package generics

import (
	"context"

	"github.com/Daniel-Njaramba-1/pulse/internal/util/interfaces"
	"github.com/jmoiron/sqlx"
)

func CreateModel[M interfaces.Store](ctx context.Context, db *sqlx.DB, model M) (int, error) {
	query := model.FeedCreateQuery()
	rows, err := db.NamedQuery(query, model)

	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var id int

	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}

		if model.FeedGetId() != nil {
			*model.FeedGetId() = id 
		}
	}
	return id, nil
}

func SelectModelById[M interfaces.Store](ctx context.Context, db *sqlx.DB, id int, model M) error {
	query := model.FeedGetByIdQuery()
	return db.Get(model, query, id)
}

func SelectAllModels[M interfaces.Store](ctx context.Context, db *sqlx.DB, models *[]M) error {
	var model M 
	query := model.FeedGetAllQuery()
	return db.Select(models, query)
}

func UpdateModelDetails[M interfaces.Store](ctx context.Context, db *sqlx.DB, model M) error {
	query := model.FeedUpdateDetailsQuery()
	_, err := db.NamedExec(query, model)
	return err
}

func DeactivateModel[M interfaces.Store](ctx context.Context, db *sqlx.DB, model M) error {
	query := model.FeedDeactivateQuery()
	_, err := db.NamedExec(query, model)
	return err
}

func ReactivateModel[M interfaces.Store](ctx context.Context, db *sqlx.DB, model M) error {
	query := model.FeedReactivateQuery()
	_, err := db.NamedExec(query, model)
	return err
}

func DeleteModel[M interfaces.Store](ctx context.Context, db *sqlx.DB, model M) error {
	query := model.FeedDeleteQuery()
	_, err := db.NamedExec(query, model)
	return err
}