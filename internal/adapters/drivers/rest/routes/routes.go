// Package routes
package routes

import (
	"fmt"

	docs "github.com/anderson89marques/roxb3/docs"
	"github.com/anderson89marques/roxb3/internal/adapters/drivers/rest/handlers/stock"
	"github.com/anderson89marques/roxb3/internal/infra/config"
	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func swaggerInfo() {
	docs.SwaggerInfo.BasePath = "/roxb3"
	docs.SwaggerInfo.Title = "Stocks b3"
	docs.SwaggerInfo.Description = "Interface consultiva para dados agregados por ativos da b3"
	docs.SwaggerInfo.Version = "1.0.0"
}

func Register(router gin.IRouter) {
	// TODO: : swagger info
	swaggerInfo()
	basePath := fmt.Sprintf("/%s", config.GetEnv().BasePath)
	fmt.Println(basePath)
	baseGroup := router.Group(basePath)
	baseGroup.GET("/docs/*any", ginSwagger.WrapHandler(swagFiles.Handler))
	baseGroup.GET("/stocks", stock.Handler)
}
