package server

import (
	"encoding/json"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/utils"
	"log"
	"net/http"
)

func (s *Server) GetCart(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseID(r.PathValue("userId"))
	if err != nil {
		log.Printf("GetCart, parse userId: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cart, err := s.cartService.GetCart(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &model.GetCartResponse{
		Items:      cart.Items,
		TotalPrice: cart.TotalPrice,
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
