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
		X string `json:"x"`
		Y int    `json:"y"`
	}
	var r []Chart
	var cc []Chart
	db.Conn.Raw("SELECT to_char(ot.created_at,'MON') as x,SUM(ot.quantity*p.price) as y FROM order_items as ot JOIN products as p ON ot.product_id=p.id WHERE p.store_id=@Store GROUP BY to_char(ot.created_at,'MON')", StoreArgument{Store: uint(storeId)}).Scan(&r)
	for _, c := range r {
		cc = append(cc, Chart{
			X: c.X,
			Y: c.Y,
		})
	}
	type Pie struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	var pie []Pie
	var pp []Pie
	db.Conn.Raw("SELECT categories.id as id,categories.name as name ,COUNT(products.id) as value from products JOIN categories on products.category_id = categories.id JOIN stores on stores.id=products.store_id WHERE products.deleted_at is null and stores.id=@Store GROUP by categories.id", StoreArgument{Store: uint(storeId)}).Scan(&pie)
	for _, p := range pie {
		pp = append(pp, Pie{
			ID:    p.ID,
			Name:  p.Name,
			Value: p.Value,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"total_price": humanize.Commaf(float64(totalPrice)),
		"count_order": humanize.Commaf(float64(corder)),

		"count_product": humanize.Commaf(float64(cproduct)),
		"pie":           pp,
		"chart":         cc,
	})
}
