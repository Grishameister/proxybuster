package httpProxy

type IRepository interface {
	StoreRequest(url string, req string) error
}
