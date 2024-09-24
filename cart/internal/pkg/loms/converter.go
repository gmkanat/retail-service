package loms

import (
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/cart/model"
	loms "gitlab.ozon.dev/kanat_9999/homework/cart/pkg/api/proto/v1"
)

func ToOrderCreateRequest(userID int64, items []model.CartItem) *loms.OrderCreateRequest {
	reqItems := make([]*loms.Item, 0, len(items))
	for _, item := range items {
		reqItems = append(reqItems, &loms.Item{
			Sku:   uint32(item.SkuId),
			Count: uint32(item.Count),
		})
	}

	return &loms.OrderCreateRequest{
		UserId: userID,
		Info: &loms.OrderInfo{
			Items: reqItems,
		},
	}
}

func ToGetStockRequest(skuID int64) *loms.StocksInfoRequest {
	return &loms.StocksInfoRequest{Sku: uint32(skuID)}
}
