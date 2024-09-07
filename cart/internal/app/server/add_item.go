package server

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/utils"
	"net/http"
)

type AddItemRequest struct {
	Count uint16 `json:"count" validate:"required,min=1"`
}

func (s *Server) AddItem(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseID(r.PathValue("userId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	skuId, err := utils.ParseID(r.PathValue("skuId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Count < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&AddItemRequest{Count: req.Count}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.cartService.AddItem(r.Context(), userId, skuId, req.Count); err != nil {
		if errors.Is(err, customerrors.SkuNotFound) {
			w.WriteHeader(http.StatusPreconditionFailed)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
