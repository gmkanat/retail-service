package suite

import (
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/app/server"
	"net/http"
)

func setupRouter(srv *server.Server) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{userId}/cart/{skuId}", srv.AddItem)
	mux.HandleFunc("GET /user/{userId}/cart", srv.GetCart)
	mux.HandleFunc("DELETE /user/{userId}/cart/{skuId}", srv.RemoveItem)
	mux.HandleFunc("DELETE /user/{userId}/cart", srv.ClearCart)
	return mux
}
