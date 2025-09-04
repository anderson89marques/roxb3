// Package database interfces
package database

import (
	"github.com/anderson89marques/roxb3/internal/adapters/drivers/rest/handlers/schemas"
	"github.com/anderson89marques/roxb3/internal/core/domain"
)

type Repository interface {
	BulkInsert([]*domain.Stock) error
	Search(*schemas.StockSchema) (*domain.StockSummary, error)
	RefreshStockSummaryMeterializedView() error
	SaveStockFile(name string) error
	IsProcessed(name string) (bool, error)
}
