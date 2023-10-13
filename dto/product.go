package dto

import "RmuttPlace/model"

type ProductRead struct {
	ID        uint
	Name      string
	Desc      string
	Category  model.CategoryRead
	Available int
	Price     int
	Weight    int
	Image     string
	Rating    float32
}

type ProductReadOne struct {
	ID        uint
	Name      string
	Desc      string
	Category  model.CategoryRead
	Store     StoreRead
	Image     string
	Available int
	Price     int
	Weight    int
	Reviews   []ReviewBodyRead
	Rating    float32
}
type ReviewBodyRead struct {
	UserID  int
	Comment string
	Rating  float32
	User    UserReview
}
type UserReview struct {
	ID    uint
	Name  string
	Image string
}
