package server

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/utils"
	"log"
	"net/http"
)

type AddItemRequest struct {
	Count uint16 `json:"count" validate:"required,min=1"`
}

func (s *Server) AddItem(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseID(r.PathValue("userId"))
	if err != nil {
		log.Printf("AddItem, parse userId: %v", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	skuId, err := utils.ParseID(r.PathValue("skuId"))
	if err != nil {
		log.Printf("AddItem, parse skuId: %v", err)
		http.Error(w, "Invalid SKU ID", http.StatusBadRequest)
		return
	}

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Count < 1 {
		log.Printf("AddItem, decode request body: %v", err)
		http.Error(w, "Invalid request body or count must be at least 1", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&AddItemRequest{Count: req.Count}); err != nil {
		log.Printf("AddItem, validation error: %v", err)
		http.Error(w, "Validation error: count must be at least 1", http.StatusBadRequest)
		return
	}

	if err := s.cartService.AddItem(r.Context(), userId, skuId, req.Count); err != nil {
		switch {
		case errors.Is(err, customerrors.SkuNotFound):
			log.Printf("AddItem, SKU not found: %v", err)
			http.Error(w, "SKU not found", http.StatusPreconditionFailed)
		case errors.Is(err, customerrors.NotEnoughStock):
			log.Printf("AddItem, not enough stock: %v", err)
			http.Error(w, "Not enough stock", http.StatusPreconditionFailed)
		default:
			log.Printf("AddItem, internal server error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
