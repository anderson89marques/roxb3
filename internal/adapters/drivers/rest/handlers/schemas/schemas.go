// Package schemas
package schemas

import "time"

type StockSchema struct {
	Ticker   string     `form:"ticker" binding:"required"`
	StarDate *time.Time `form:"data_inicio" time_format:"2006-01-02"`
}
