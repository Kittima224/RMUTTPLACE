package dto

type FavoriteRequest struct {
	UserID    uint `json:"userId"`
	ProductID uint `json:"productId"`
}
type FavoriteReponse struct {
	ProductID uint `json:"productId"`
}
