package model

import "gorm.io/gorm"

type Shipment struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;type:varchar(100);not null"`
}
type ShipmentRead struct {
	ID   uint
	Name string
}
