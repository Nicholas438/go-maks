package request

type CoinCreateRequest struct {
	Name string `json:"name" validate:"required"`
}
