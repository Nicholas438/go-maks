package request

type TradeCreateRequest struct {
	Price  int `json:"price" validate:"required"`
	CoinID int `json:"coin_id" validate:"required"`
}
