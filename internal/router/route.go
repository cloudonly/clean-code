package router

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/luomu/clean-code/gen/apis/luomu/greet/v1/greetv1connect"
	"github.com/luomu/clean-code/internal/service/greet"
	"github.com/luomu/clean-code/pkg/middleware"
)

func NewGinEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	//engine.Use(gin.Logger())
	engine.Use(middleware.Log())
	engine.Use(gin.Recovery())
	engine.Use(otelgin.Middleware("gin-server"))
	engine.UseH2C = true

	v1 := engine.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.String(200, "<h1>Hi, Clean Code</h1>")
		})
	}

	webhook := engine.Group("hooks")
	registryWebhook(webhook)

	registerConnect(engine)
	return engine
}

func registerConnect(router *gin.Engine) {

	greeter := &greet.GreetServer{}
	greetPath, greetHandler := greetv1connect.NewGreetServiceHandler(greeter)
	router.Any(greetPath+"/*w", gin.WrapH(greetHandler))
}
