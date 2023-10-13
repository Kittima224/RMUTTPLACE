package dto

type OrderRead struct {
	OrderID    uint
	ShipmentID uint
}

type OrderProductRead struct {
	ID    uint
	Name  string
	Price int
	Image string
}
type OrderRequest struct {
	StoreID uint
	Carts   []OrderItemRequest
}
type OrderItemRequest struct {
	ProductID uint
	Quantity  int
	Product   OrderProductRead
	Store     StoreRead
}

type OrderOneRead struct {
	OrderID      uint
	ShipmentID   uint
	ShipmentName string
	StoreID      uint
	Store        StoreRead
	Product      OrderProductRead
	TotalPrice   int
}

type OrderReadOne struct {
	ID       uint
	Store    StoreRead
	Products []OrderItemRead
	Shipment ShipmentRead
	Tracking string
	User     UserAddress
}
type OrderItemRead struct {
	ID       uint
	Name     string
	Price    int
	Image    string
	Quantity int
}
type OrderReadAll struct {
	ID           uint
	UserID       uint
	ShipmentID   uint
	ShipmentName string
	Tracking     string
	Store        StoreRead
}
