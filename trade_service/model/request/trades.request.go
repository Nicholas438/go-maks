package request

type TradeCreateRequest struct {
	Price  float64 `json:"price" validate:"required"`
	CoinID int     `json:"coin_id" validate:"required,gt=0"`
}
