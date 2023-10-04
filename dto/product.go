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
	Rating    int
}

type ProductReadOne struct {
	ID        uint
	Name      string
	Desc      string
	Category  model.CategoryRead
	Store     model.StoreRead
	Image     string
	Available int
	Price     int
	Weight    int
	Reviews   []ReviewBodyRead
	Rating    int
}
type ReviewBodyRead struct {
	UserID  int
	Name    string
	Comment string
	Rating  int
}
