package main

import (
	AdminController "RmuttPlace/controller/admin"
	AuthController "RmuttPlace/controller/auth"
	StoreController "RmuttPlace/controller/store"
	UntokenController "RmuttPlace/controller/untoken"
	UserController "RmuttPlace/controller/user"
	"RmuttPlace/db"
	"RmuttPlace/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	// err := godotenv.Load(".env")

	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// }
	db.InitDB()
	db.Migrate()

	r := gin.Default()
	r = gin.New()
	r.Use(CORSMiddleware())
	os.MkdirAll("uploads/products", 0755)
	os.MkdirAll("uploads/stores", 0755)
	os.MkdirAll("uploads/admins", 0755)
	os.MkdirAll("uploads/users", 0755)
	os.MkdirAll("uploads/files", 0755)
	r.Static("/uploads", "./uploads")
	r.POST("/user/register", AuthController.Register)
	r.POST("/user/login", AuthController.Login)
	r.POST("/store/register", AuthController.RegisterStore)
	r.POST("/store/login", AuthController.LoginStore)
	r.POST("/admin/register", AuthController.RegisterAdmin)
	r.POST("/admin/login", AuthController.LoginAdmin)

	admin := r.Group("/admins", middleware.JWTAuthenAdmin())
	admin.GET("/profile", AdminController.Profile)
	admin.GET("/adminprofile", AdminController.GetProfile)
	admin.PATCH("/adminprofile", AdminController.UpdateAdmin)
	admin.GET("/all", AdminController.ReadAll)

	admin.POST("/category", AdminController.Create)
	admin.GET("/categories", AdminController.CategoryAll)
	admin.GET("/category/:id", AdminController.CategoryOne)
	admin.DELETE("/category", AdminController.CategoryDel)
	admin.PATCH("/category/:id", AdminController.CategoryUpdate)

	admin.GET("/products", AdminController.ReadProductAll)
	admin.DELETE("/product/del", AdminController.DeleteProduct)
	admin.GET("/product/:id", AdminController.ReadOneProduct)
	admin.PATCH("/product/update/:id", AdminController.UpdateProduct)

	admin.GET("/dashboard", AdminController.Dashboard)

	admin.GET("/orders", AdminController.GetOrderAll)
	admin.GET("/order/:id", AdminController.GetOrderOne)

	admin.POST("/store", AdminController.StoreRegister)
	admin.GET("/stores", AdminController.ReadAllStore)
	admin.GET("/store/:id", AdminController.ReadOneStore)
	admin.PATCH("/store/update/:id", AdminController.UpdateStore)
	admin.DELETE("/store/del", AdminController.DeleteStore)

	admin.POST("/shipment", AdminController.ShipmentCreate)
	admin.GET("/shipments", AdminController.ShipmentAll)
	admin.GET("/shipment/:id", AdminController.ShipmentOne)
	admin.PATCH("/shipment/:id", AdminController.ShipmentUpdate)
	admin.DELETE("/shipment", AdminController.ShipmentDel)

	admin.POST("/user", AdminController.Register)
	admin.GET("/users", AdminController.ReadAllUser)
	admin.GET("/user/:id", AdminController.ReadOneUser)
	admin.PATCH("/user/update/:id", AdminController.UpdateUser)
	admin.DELETE("/user/del", AdminController.DeleteUser)

	user := r.Group("/users", middleware.JWTAuthen())
	user.GET("/readall", UserController.ReadAll)
	user.GET("/profile", UserController.Profile)
	user.PATCH("/update", UserController.UpdateProfileUser)

	user.POST("/mycart", UserController.AddCart)
	user.GET("/mycart", UserController.MyCart)
	user.DELETE("/mycart/:id", UserController.DeleteProductMyCart)
	user.PATCH("/mycart/quantity", UserController.UpdateQuantity)

	user.POST("/checkout", UserController.CreateOrder)
	// user.POST("/buynow", UserController.BuyNow)
	user.GET("/myorder", UserController.MyOrderAll)
	user.GET("/myorder/:id", UserController.MyOrderFindOne)

	user.POST("/fav", UserController.AddFav)
	user.DELETE("/fav", UserController.UnFav)
	user.GET("/fav", UserController.MyFav)

	user.POST("/product/review/:id", UserController.CreateReview)

	store := r.Group("/stores", middleware.JWTAuthenStore())
	store.GET("/readall", StoreController.ReadAll)
	store.GET("/profile", StoreController.Profile)
	store.PATCH("/profile", StoreController.UpdateMyStore)
	// store.PATCH("/changeprofile", StoreController.ChangProfileMystore) chang pass

	store.POST("/product", StoreController.Create)
	store.GET("/products", StoreController.ReadProductAllMyStore) //รวมกับ findnameproduct
	store.GET("/product/:id", StoreController.FindOneProductMyStore)
	store.PATCH("/product/:id", StoreController.UpdateProductMystore)
	store.DELETE("/product/:id", StoreController.DeleteProductMyStore)

	store.GET("/orders", StoreController.GetOrderAll)
	store.GET("/order/:id", StoreController.GetOrderOne)
	store.PATCH("/order/:id", StoreController.AddTrackingOrder)

	store.GET("/dashboard", StoreController.DashboardStore)
	//untoken
	r.GET("/store/:id/products", UntokenController.ProductAllStore)
	r.GET("/products", UntokenController.ReadProductAll) //ค้นชื่อสินค้า ชื่อหมสดหมู่categoty ค้นแท็กในdesc ได้
	r.GET("/product/:id", UntokenController.FindOneProduct)
	r.GET("/store/:id", UntokenController.ReadOneStore)
	//admin //admin //admin //admin //admin //admin //admin

	//test
	r.GET("/dashboard", AdminController.DashboardTest)

	//http.ListenAndServe(":3000", nil)
	//r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	port := os.Getenv("PORT")
	if port != "" {
		port = "8080"
	}
	r.Run(":" + port)
}
