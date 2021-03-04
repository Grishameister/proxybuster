package delivery

import (
	"bufio"
	"github.com/Grishameister/proxybuster/configs/config"
	"github.com/Grishameister/proxybuster/pkg/api"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type ApiHandler struct {
	client http.Client
	repo   api.IRepository
}

func NewApiHandler(client http.Client, repo api.IRepository) *ApiHandler {
	return &ApiHandler{
		client: client,
		repo:   repo,
	}
}

func (h *ApiHandler) newReq(c *gin.Context) (*http.Request, string, error) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, "", err
	}

	buffReq, err := h.repo.GetRequest(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, "", err
	}

	b := bufio.NewReader(strings.NewReader(buffReq.Req))

	req, err := http.ReadRequest(b)
	if err != nil {
		config.Lg("repeat", "createReq").Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, "", err
	}

	req.RequestURI = ""
	return req, buffReq.Url, nil
}

func (h *ApiHandler) HandleRepeat(c *gin.Context) {
	req, urlReq, err := h.newReq(c)
	if err != nil {
		return
	}

	u, err := url.Parse(urlReq)
	if err != nil {
		config.Lg("repeat", "url").Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	req.URL = u

	resp, err := h.client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		config.Lg("repeat", "Do").Error(err.Error())
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, v := range values {
			c.Writer.Header().Add(key, v)
		}
	}

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		config.Lg("repeat", "Copy").Error(err.Error())
		return
	}
}

func (h *ApiHandler) ScanHandler(c *gin.Context) {
	req, urlReq, err := h.newReq(c)
	if err != nil {
		return
	}

	i := 0
	var b strings.Builder
	idx := len(urlReq)
	for j, c := range urlReq {
		if c == '/' {
			i++
		}
		if i == 3 {
			idx = j
			break
		}
	}

	file, err := os.Open("dicc.txt")

	if err != nil {
		config.Lg("buster", "file").Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(file)

	var files []string
	for scanner.Scan() {
		b.WriteString(urlReq[:idx])
		b.WriteByte('/')
		qs := url.QueryEscape(scanner.Text())
		b.WriteString(qs)

		u, err := url.Parse(b.String())
		if err != nil {
			config.Lg("repeat", "url").Error(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		req.URL = u
		resp, err := h.client.Do(req)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			config.Lg("repeat", "Do").Error(err.Error())
			return
		}

		if resp.StatusCode != http.StatusBadRequest {
			files = append(files, scanner.Text())
		}
		b.Reset()
	}

	c.JSON(200, files)
}
