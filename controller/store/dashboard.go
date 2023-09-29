package store

import (
	"RmuttPlace/db"
	"net/http"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
)

func DashboardStore(c *gin.Context) {
	storeId := c.MustGet("storeId").(float64)
	type StoreArgument struct {
		Store uint
	}
	var totalPrice int
	db.Conn.Raw("SELECT sum(order_items.quantity*products.price) from order_items JOIN products ON products.id = order_items.product_id JOIN stores ON products.store_id=stores.id WHERE stores.id = @Store",
		StoreArgument{Store: uint(storeId)}).Scan(&totalPrice)

	var cproduct int
	db.Conn.Raw("SELECT COUNT(products.id) from products JOIN stores ON stores.id=products.store_id WHERE products.deleted_at is null and stores.id= @Store", StoreArgument{Store: uint(storeId)}).Scan(&cproduct)

	var corder int
	db.Conn.Raw("SELECT COUNT(id) from orders WHERE deleted_at is null and store_id= @Store", StoreArgument{Store: uint(storeId)}).Scan(&corder)

	type Chart struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	var r []Chart
	db.Conn.Raw("SELECT DATE_FORMAT(ot.created_at,'%M') as name,SUM(ot.quantity*p.price) as value,DATE_FORMAT(ot.created_at,'%m') as id FROM `order_items` as ot JOIN `products` as p ON ot.product_id=p.id WHERE p.store_id=@Store GROUP BY DATE_FORMAT(created_at,'%M')", StoreArgument{Store: uint(storeId)}).Scan(&r)
	c.JSON(http.StatusOK, gin.H{
		"total_price": humanize.Commaf(float64(totalPrice)),
		"count_order": humanize.Commaf(float64(corder)),
		// "count_acc":     humanize.Commaf(float64(cuser)),
		"count_product": humanize.Commaf(float64(cproduct)),
		// "pie":           result,
		"chart": r,
	})
}
