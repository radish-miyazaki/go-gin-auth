package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/radish-miyazaki/go-auth/controllers"
)

func Setup(r *gin.Engine) {
	// api/v1
	v1 := r.Group("api/v1")
	{
		// common
		v1.POST("/register", controllers.Register)
		v1.POST("/login", controllers.Login)
		v1.POST("/logout", controllers.Logout)
		v1.GET("/user", controllers.User)
		v1.POST("/forgot", controllers.Forgot)
		v1.POST("/reset", controllers.Reset)
	}
}
