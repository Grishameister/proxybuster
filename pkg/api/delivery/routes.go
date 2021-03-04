package delivery

import (
	"github.com/Grishameister/proxybuster/internal/database"
	"github.com/Grishameister/proxybuster/pkg/api/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddRoutes(r *gin.Engine, db database.IDbConn) {
	repo := repository.NewRepoApi(db)
	handler := NewApiHandler(http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, repo)

	r.GET("/requests", handler.GetRequests)
	r.GET("/requests/:id", handler.GetRequest)

	r.GET("/repeat/:id", handler.HandleRepeat)
	r.GET("/scan/:id", handler.ScanHandler)
}
