package main

import (
	"go_jwt/controllers"
	"go_jwt/initializers"

	"go_jwt/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariable()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("Templates/*.html")
	router.Use(middleware.RequireAuth)

	router.GET("/signup", controllers.Loadsignup)
	router.POST("/signup", controllers.SignUp)
	router.GET("/success", controllers.LoadSuccess)
	router.GET("/login", controllers.LoadLogin)
	router.GET("/home", controllers.Loadhome)
	router.POST("/login", controllers.Login)
	router.GET("/admin", controllers.AdminPageHandler)
	router.POST("/admin/users/:id/delete", controllers.DeleteUser)
	router.POST("/admin/users/:id/update", controllers.UpdateUser)
	router.GET("/back", controllers.Signout)
	//router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	router.Run()

}
