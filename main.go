package main

import (
	AdminController "RmuttPlace/controller/admin"
	AuthController "RmuttPlace/controller/auth"
	StoreController "RmuttPlace/controller/store"
	UserController "RmuttPlace/controller/user"
	"RmuttPlace/db"
	"RmuttPlace/middleware"
	"fmt"

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
	// if os.Getenv("APP_ENV") == "production" {
	// 	gin.SetMode(gin.ReleaseMode)
	// } else {
	// 	if err := godotenv.Load(); err != nil {
	// 		log.Fatal("Error loading .env file")
	// 	}
	// }
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}
	db.InitDB()
	db.Migrate()

	r := gin.Default()
	r = gin.New()
	r.Use(CORSMiddleware())
	r.Static("/uploads", "./uploads")
	r.POST("/user/register", AuthController.Register)
	r.POST("/user/login", AuthController.Login)
	r.POST("/store/register", AuthController.RegisterStore)
	r.POST("/store/login", AuthController.LoginStore)
	r.POST("/admin/register", AuthController.RegisterAdmin)
	r.POST("/admin/login", AuthController.LoginAdmin)

	r.GET("/product/findname", StoreController.FindNameProduct)

	admin := r.Group("/admins", middleware.JWTAuthenAdmin())
	admin.GET("/profile", AdminController.Profile)
	admin.GET("/all", AdminController.ReadAll)

	user := r.Group("/users", middleware.JWTAuthen())
	user.GET("/readall", UserController.ReadAll)
	user.GET("/profile", UserController.Profile)
	user.PATCH("/addinfo", UserController.AddProfileUser)
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

	//admin //admin //admin //admin //admin //admin //admin
	r.POST("/admin/store", AdminController.StoreRegister)
	r.GET("/admin/stores", AdminController.ReadAllStore)
	r.GET("/admin/store/:id", AdminController.ReadOneStore)
	r.PATCH("/admin/store/update/:id", AdminController.UpdateStore)
	r.DELETE("/admin/store/del", AdminController.DeleteStore)

	r.GET("/admin/products", AdminController.ReadProductAll)
	r.DELETE("/admin/product/del", AdminController.DeleteProduct)
	r.GET("/admin/product/:id", AdminController.ReadOneProduct)
	r.PATCH("/admin/product/update/:id", AdminController.UpdateProduct)

	r.POST("/admin/user", AdminController.Register)
	r.GET("/admin/users", AdminController.ReadAllUser)
	r.GET("/admin/user/:id", AdminController.ReadOneUser)
	r.PATCH("/admin/user/update/:id", AdminController.UpdateUser)
	r.PATCH("/admin/user/update/:id/photo", AdminController.UpdateUserPhoto)
	r.DELETE("/admin/user/del", AdminController.DeleteUser)

	r.GET("/admin/categories", AdminController.CategoryAll)
	r.GET("/admin/category/:id", AdminController.CategoryOne)
	r.DELETE("/admin/category", AdminController.CategoryDel)
	r.PATCH("/admin/category/:id", AdminController.CategoryUpdate)
	r.POST("/admin/category", AdminController.Create)

	r.POST("/admin/shipment", AdminController.ShipmentCreate)
	r.GET("/admin/shipments", AdminController.ShipmentAll)
	r.GET("/admin/shipment/:id", AdminController.ShipmentOne)
	r.PATCH("/admin/shipment/:id", AdminController.ShipmentUpdate)
	r.DELETE("/admin/shipment", AdminController.ShipmentDel)

	r.GET("/admin/orders", AdminController.GetOrderAll)
	r.GET("/admin/order/:id", AdminController.GetOrderOne)

	r.GET("/admin/dashboard", AdminController.Dashboard)
	//test

	//http.ListenAndServe(":3000", nil)
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// port := os.Getenv("PORT")
	// if port != "" {
	// 	port = "8080"
	// }
	// r.Run(":" + port)
}
