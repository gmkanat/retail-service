package service

import (
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/http/transport"
	"net/http"
)

type Product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type ProductService struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewProductService(baseURL, token string) *ProductService {
	return &ProductService{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Transport: &transport.RetryRoundTripper{
				Next: http.DefaultTransport,
			},
		},
	}
}
