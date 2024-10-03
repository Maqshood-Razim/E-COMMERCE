package routes

import (
	"furniture-ecommerce/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.Default()

	// User routes
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/products", controllers.BrowseProducts)
	r.GET("/products/search", controllers.SearchProduct)

	// Protected routes (require JWT authentication)
	protected := r.Group("/")
	protected.Use(controllers.JWTAuthMiddleware())
	{
		protected.POST("/cart/add", controllers.AddToCart)
		protected.DELETE("/cart/remove/:id", controllers.RemoveFromCart)
		protected.GET("/cart", controllers.GetCart)
		protected.GET("/order/details", controllers.OrderDetails)
		protected.POST("/payment/process", controllers.ProcessPayment)
		protected.POST("/checkout", controllers.Checkout)
	}

	// Admin login route (no middleware needed)
	r.POST("/admin/login", controllers.AdminLogin)

	// Initialize admin routes
	AdminRoutes(r)  

	return r
}

func AdminRoutes(router *gin.Engine) {
	
	adminGroup := router.Group("/admin")
	adminGroup.Use(controllers.JWTAuthMiddleware(), controllers.AdminOnly()) // Apply JWT and admin check to the rest of the routes
	{
		adminGroup.GET("/users", controllers.ListUsers)
		adminGroup.POST("/addproducts", controllers.AddProduct) 
		adminGroup.PUT("/products/:id", controllers.EditProduct)
		adminGroup.DELETE("/products/:id", controllers.DeleteProduct)
		adminGroup.DELETE("/users/:id", controllers.DeleteUser)
	}
}
