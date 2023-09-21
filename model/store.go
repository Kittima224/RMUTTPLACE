package model

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	UserName    string
	Tel         string
	Email       string
	NameStore   string
	Address     string
	District    string
	SubDistrict string
	Province    string
	Zipcode     string
	Password    string
	Image       string
	// Shipment      string
	// Shipments     []Shipment `gorm:"foreignKey:StoreID"`
	Status        bool //ยังไม่ใส่ค่า
	AccountNumber string
	AccountName   string
	Bank          string
}

//`gorm:"foreignKey:StoreID"`

// type Shipment struct {
// 	StoreID      uint
// 	ShipmentName string
// }
