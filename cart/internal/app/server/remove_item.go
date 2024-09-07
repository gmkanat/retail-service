package server

import (
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/utils"
	"log"
	"net/http"
)

func (s *Server) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ParseID(r.PathValue("userId"))
	if err != nil {
		log.Printf("RemoveItem, parse userId: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	skuId, err := utils.ParseID(r.PathValue("skuId"))
	if err != nil {
		log.Printf("RemoveItem, parse skuId: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.cartService.RemoveItem(r.Context(), userId, skuId); err != nil {
		log.Printf("RemoveItem: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
