// Package domain
package domain

import (
	"strconv"
	"strings"
	"time"

	"github.com/anderson89marques/roxb3/internal/core/utils"
)

type Stock struct {
	Ticker      string    // CodigoInstrumento
	TradeDate   time.Time // DataNegocio
	TradePrice  float64   // PrecoNegocio
	Volume      int64     // QuantidadeNegociada
	ClosingTime time.Time // HoraFechamento
}

type StockSummary struct {
	Ticker         string  `json:"ticker" binding:"required"`
	MaxRangeValue  float64 `json:"max_range_value" binding:"required"`
	MaxDailyVolume float64 `json:"max_daily_volume" binding:"required"`
}

const dateFormat = "2006-01-02"

func ToStockModel(row []string) (*Stock, error) {
	ticker := row[1]
	tradeDate, err := time.Parse(dateFormat, row[8])
	if err != nil {
		return nil, err
	}
	periodPrice := strings.Replace(row[3], ",", ".", 1)
	tradePrice, err := strconv.ParseFloat(periodPrice, 64)
	if err != nil {
		return nil, err
	}
	volume, err := strconv.Atoi(row[4])
	if err != nil {
		return nil, err
	}

	closingTime, err := utils.ParseStringHourMilisToTime(row[5])
	if err != nil {
		return nil, err
	}

	return &Stock{
		Ticker:      ticker,
		TradeDate:   tradeDate,
		TradePrice:  tradePrice,
		Volume:      int64(volume),
		ClosingTime: *closingTime,
	}, nil
}
