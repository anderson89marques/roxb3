// Package Stock
package stock

import (
	"net/http"

	database "github.com/anderson89marques/roxb3/internal/adapters/drivens/postgres"
	"github.com/anderson89marques/roxb3/internal/adapters/drivers/rest/handlers/schemas"
	"github.com/anderson89marques/roxb3/internal/infra/config"
	"github.com/anderson89marques/roxb3/internal/services"
	"github.com/gin-gonic/gin"
)

// @Summary Consulta de dados agregados por ativo
// @Description A consulta retorna um json com o maior preco unitário e volume máximo de ativo negociado por ativoco unitário e volume máximo de ativo negociado por ativo
// @Tags Stock
// @Accept json
// @Produce json
// @Param ticker query string true "ativo"
// @Param data_inicio query string false "Data negócio no formato yyyy-mm-dd"
// @Success 200 {object} domain.StockSummary
// @Router /stocks [get]
func Handler(c *gin.Context) {
	repo, err := database.NewRepository(config.GetEnv())
	if err != nil {
		panic(err)
	}
	stockService := services.NewStockService(repo, config.GetEnv())
	var stockParams schemas.StockSchema
	if err := c.ShouldBindQuery(&stockParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stockSummary, err := stockService.Search(&stockParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stockSummary)
}
