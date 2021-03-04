package repository

import (
	"context"
	"github.com/Grishameister/proxybuster/configs/config"
	"github.com/Grishameister/proxybuster/internal/database"
)

type Repository struct {
	db database.IDbConn
}

func NewRepo(db database.IDbConn) *Repository {
	return &Repository{
		db:db,
	}
}

func (r *Repository) StoreRequest(url string, req string) error {
	if _, err := r.db.Exec(context.Background(), "insert into requests(url, body) values($1, $2)", url, req); err != nil {
		config.Lg("repository", "StoreRequest").Error(err.Error())
		return err
	}
	return nil
}
