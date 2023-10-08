package dto

type ShipmentBody struct {
	Name string `json:"name" binding:"required"`
}
