package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type CheckoutRequest struct {
	UserId int64 `json:"user"`
}

type CheckoutResponse struct {
	OrderId int64 `json:"orderID"`
}

func (s *Server) Checkout(w http.ResponseWriter, r *http.Request) {
	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Checkout, decode request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	orderId, err := s.cartService.Checkout(r.Context(), req.UserId)
	if err != nil {
		log.Printf("Checkout, service error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := CheckoutResponse{OrderId: orderId}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Checkout, encode response: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
