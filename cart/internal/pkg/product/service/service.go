package service

import (
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

func NewProductService(baseURL, token string, client *http.Client) *ProductService {
	return &ProductService{
		baseURL: baseURL,
		token:   token,
		client:  client,
	}
}
