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
		log.Printf("GetCart, parse userId: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cart, err := s.cartService.GetCart(r.Context(), userId)
	if err != nil {
		log.Printf("GetCart, fetch cart: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(cart.Items) == 0 {
		http.Error(w, "Cart is empty", http.StatusNotFound)
		return
	}

	response := &model.GetCartResponse{
		Items:      cart.Items,
		TotalPrice: cart.TotalPrice,
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("GetCart, encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
