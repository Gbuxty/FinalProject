package router

import (
	"gw-currency-wallet/internal/transport/http/handlers"
	"gw-currency-wallet/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Routers struct {
	router           *gin.Engine
	logger           logger.Logger
	authHandler      *handlers.AuthHandler
	walletHandler    *handlers.WalletHandlers
	exchangerHangler *handlers.ExchangeHandlers
}

func New(
	logger logger.Logger,
	authHandler *handlers.AuthHandler,
	walletHandler *handlers.WalletHandlers,
	exchangerHangler *handlers.ExchangeHandlers,
) *Routers {
	return &Routers{
		router:           gin.New(),
		logger:           logger,
		authHandler:      authHandler,
		walletHandler:    walletHandler,
		exchangerHangler: exchangerHangler,
	}
}

func (r *Routers) registerAuthRoutes(group *gin.RouterGroup) {
	group.POST("/register", r.authHandler.Register)
	group.POST("/login", r.authHandler.Login)
}

func (r *Routers) registerWalletRoutes(group *gin.RouterGroup) {
	group.Use(r.authHandler.AuthRequired())
	group.GET("/balance", r.walletHandler.GetBalance)
	group.POST("/deposit", r.walletHandler.Deposit)
	group.POST("/withdraw", r.walletHandler.Withdraw)
}

func (r *Routers) registerExchangertRoutes(group *gin.RouterGroup) {
	group.Use(r.authHandler.AuthRequired())
	group.GET("/rates", r.exchangerHangler.GetAllRates)
	group.POST("/excurrency",r.exchangerHangler.ChangeCurrency)
}

func (r *Routers) Routes() *gin.Engine {
	apiV1 := r.router.Group("api/v1")

	r.registerAuthRoutes(apiV1.Group(""))
	r.registerWalletRoutes(apiV1.Group("wallet"))
	r.registerExchangertRoutes(apiV1.Group("exchange"))

	return r.router
}
