package domain

type Request struct {
	Req string `json:"req"`
	Url string `json:"url"`
	Id  int    `json:"id"`
}
