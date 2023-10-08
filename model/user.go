package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string
	UserName    string
	Password    string
	Address     string
	District    string
	SubDistrict string
	Province    string
	Zipcode     string
	Tel         string
	Image       string
}




