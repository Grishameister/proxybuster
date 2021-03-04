package repository

import (
	"context"
	"github.com/Grishameister/proxybuster/configs/config"
	"github.com/Grishameister/proxybuster/internal/database"
	"github.com/Grishameister/proxybuster/pkg/domain"
)

type Repository struct {
	db database.IDbConn
}

func NewRepoApi(db database.IDbConn) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetRequests() ([]domain.Request, error) {
	rows, err := r.db.Query(context.Background(), "select id, url, body from requests")
	if err != nil {
		config.Lg("apiRepo", "GetRequests").Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var requests []domain.Request
	var req domain.Request

	for rows.Next() {
		if err := rows.Scan(&req.Id, &req.Url, &req.Req); err != nil {
			config.Lg("apiRepo", "GetRequestsScan").Error(err.Error())
			return nil, err
		}

		requests = append(requests, req)
	}

	return requests, nil
}

func (r *Repository) GetRequest(id int) (domain.Request, error) {
	d := domain.Request{
		Id: id,
	}

	if err := r.db.QueryRow(context.Background(), "select url, body from requests where id = $1", d.Id).Scan(&d.Url, &d.Req); err != nil {
		config.Lg("apiRepo", "GetRawReq").Error(err.Error())
		return d, err
	}
	return d, nil
}
