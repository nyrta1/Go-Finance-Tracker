package routers

import (
	"github.com/gin-gonic/gin"
	"go-finance-tracker/internal/rest/handler"
	"go-finance-tracker/pkg/middleware"
)

type Routers struct {
	authHandler    *handler.AuthHandlers
	financeHandler *handler.FinanceHandlers
}

func NewRouters(authHandler *handler.AuthHandlers, financeHandler *handler.FinanceHandlers) *Routers {
	return &Routers{
		authHandler:    authHandler,
		financeHandler: financeHandler,
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
		financeRouter := v1Router.Group("/finance", middleware.RequireAuthMiddleware)
		{
			financeRouter.GET("", r.financeHandler.GetAllFinance)
			financeRouter.POST("", r.financeHandler.AddFinanceRecord)
		}
	}
}
