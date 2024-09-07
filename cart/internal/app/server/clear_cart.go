package server

import (
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/utils"
	"log"
	"net/http"
)

func (s *Server) ClearCart(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseID(r.PathValue("userId"))
	if err != nil {
		log.Printf("ClearCart, parse userId: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.cartService.ClearCart(r.Context(), userId); err != nil {
		log.Printf("ClearCart: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
