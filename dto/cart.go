package dto

type CartResponse struct {
	UserID    uint `json:"userId"`
	ProductID uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}
type CartRequest struct {
	StoreID   uint `json:"storeId"`
	ProductID uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}
type ReadProductInCart struct {
	Store    StoreRead
	ID       uint
	Image    string
	Name     string
	Quantity int
}
