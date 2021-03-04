package delivery

import (
	"github.com/Grishameister/proxybuster/configs/config"
	"github.com/Grishameister/proxybuster/pkg/httpProxy"
	"io"
	"net/http"
	"net/http/httputil"
)

type HttpProxyHandler struct {
	client http.Client
	repo   httpProxy.IRepository
}

func NewProxyHttp(client http.Client, repo httpProxy.IRepository) *HttpProxyHandler {
	return &HttpProxyHandler{
		client: client,
		repo:   repo,
	}
}

func (h *HttpProxyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	reqNew, err := http.NewRequest(r.Method, r.RequestURI, r.Body)
	defer r.Body.Close()
	if err != nil {
		config.Lg("httpProxyHandler", "Handle").Error(err.Error())
		w.WriteHeader(400)
		return
	}

	for k, headers := range r.Header {
		if k == "Proxy-Connection" {
			continue
		}
		for _, h := range headers {
			reqNew.Header.Add(k, h)
		}
	}

	req, err := httputil.DumpRequest(reqNew, true)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if err := h.repo.StoreRequest(r.RequestURI, string(req)); err != nil {
		return
	}

	resp, err := h.client.Do(reqNew)
	if err != nil {
		config.Lg("httpProxyHandler", "Handle").Error(err.Error())
		return
	}

	w.WriteHeader(resp.StatusCode)
	for k, headers := range resp.Header {
		for _, h := range headers {
			w.Header().Add(k, h)
		}
	}

	_, err = io.Copy(w, resp.Body)
	resp.Body.Close()
}
