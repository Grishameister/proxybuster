package api

import "github.com/Grishameister/proxybuster/pkg/domain"

type IRepository interface {
	GetRequests() ([]domain.Request, error)
	GetRequest(id int) (domain.Request, error)
}
