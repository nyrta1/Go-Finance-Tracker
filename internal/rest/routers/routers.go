package routers

import (
	"github.com/gin-gonic/gin"
	"go-finance-tracker/internal/rest/handler"
	"go-finance-tracker/pkg/middleware"
)

type Routers struct {
	authHandler handler.AuthHandlers
}

func NewRouters(authHandler handler.AuthHandlers) *Routers {
	return &Routers{
		authHandler: authHandler,
	}
}

func (r *Routers) SetupRoutes(app *gin.Engine) {
	v1Router := app.Group("/v1")
	{
		authRouter := v1Router.Group("/auth")
		{
			authRouter.POST("/register", r.authHandler.Register)
			authRouter.POST("/login", r.authHandler.Login)
			authRouter.POST("/logout", middleware.RequireAuthMiddleware, r.authHandler.Logout)
			authRouter.GET("/profile", middleware.RequireAuthMiddleware, r.authHandler.Profile)
		}
	}
}
