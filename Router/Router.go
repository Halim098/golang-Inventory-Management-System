package Router

import (
	"ims/Controller"
	middleware "ims/Middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // Menambahkan "Authorization" ke AllowHeaders
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/register", Controller.Register)
	r.POST("/login", Controller.Login)
	r.GET("/logout", Controller.Logout)

	v1 := r.Group("/v1", middleware.AdminAuth())
	{
		v1.GET("/user", Controller.GetUser)
		v1.GET("/product", Controller.GetProducts)
		v1.GET("/product/:id", Controller.GetProductsByID)
		v1.POST("/product", Controller.AddProduct)
		v1.PUT("/product/:id", Controller.UpdateProduct)
		v1.DELETE("/product/:id", Controller.DeleteProduct)
		v1.GET("/activity", Controller.GetActivities)
		v1.GET("/activity/:id", Controller.GetActivityByID)
	}

	v2 := r.Group("/v2", middleware.UserAuth())
	{
		v2.GET("/user", Controller.GetUser)
		v2.GET("/product", Controller.GetProducts)
		v2.GET("/product/:id", Controller.GetProductsByID)
	}

	return r
}
