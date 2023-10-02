package dto

type CartResponse struct {
	UserID    uint `json:"userId"`
	ProductID uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}
type CartRequest struct {
	UserID    uint `json:"userId"`
	ProductID uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}
