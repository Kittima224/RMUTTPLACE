package dto

type ShipmentBody struct {
	Name string `json:"name" binding:"required"`
}
type ShipmentRead struct {
	ID   uint
	Name string
}
