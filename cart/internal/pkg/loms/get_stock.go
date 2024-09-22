package loms

import "context"

func (c *Client) GetStock(ctx context.Context, skuID int64) (int64, error) {
	resp, err := c.LomsClient.StocksInfo(ctx, ToGetStockRequest(skuID))

	if err != nil {
		return 0, err
	}

	return int64(resp.AvailableCount), nil
}
