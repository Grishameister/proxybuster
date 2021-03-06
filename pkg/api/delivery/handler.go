package delivery

import (
	"bufio"
	"github.com/Grishameister/proxybuster/configs/config"
	"github.com/Grishameister/proxybuster/pkg/api"
	"github.com/Grishameister/proxybuster/pkg/domain"
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

func (h *ApiHandler) copyReq(c *gin.Context) (domain.Request,  error) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return domain.Request{}, err
	}

	return h.repo.GetRequest(id)
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
	domreq, err := h.copyReq(c)
	if err != nil {
		return
	}

	i := 0
	var b strings.Builder
	idx := len(domreq.Url)
	for j, c := range domreq.Url {
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
	defer func() {
		if err := file.Close(); err != nil {
			config.Lg("Close", "File").Error(err.Error())
		}
		return
	}()
	scanner := bufio.NewScanner(file)

	var files []string
	for scanner.Scan() {
		tempBufReq := domreq
		bufReader := bufio.NewReader(strings.NewReader(tempBufReq.Req))

		req, err := http.ReadRequest(bufReader)
		if err != nil {
			config.Lg("repeat", "createReq").Error(err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		req.RequestURI = ""

		b.WriteString(domreq.Url[:idx])
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
			//config.Lg("repeat", "Do").Error(err.Error())
			b.Reset()
			continue
		}

		config.Lg("log", "req").Info(b.String())
		if resp.StatusCode != http.StatusNotFound {
			files = append(files, scanner.Text())
		}
		b.Reset()
	}

	c.JSON(200, files)
}


func (h *ApiHandler) GetRequest(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	buffReq, err := h.repo.GetRequest(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, buffReq)
}

func (h *ApiHandler) GetRequests(c *gin.Context) {
	buffReq, err := h.repo.GetRequests()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, buffReq)
}
