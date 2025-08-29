// Package database
package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/anderson89marques/roxb3/internal/adapters/drivers/rest/handlers/schemas"
	"github.com/anderson89marques/roxb3/internal/core/domain"
	"github.com/anderson89marques/roxb3/internal/infra/config"
)

type repository struct {
	db *sql.DB
}

func NewRepository(conf *config.Config) (*repository, error) {
	database, err := NewPostgresDB(conf)
	if err != nil {
		return nil, err
	}

	return &repository{
		db: database,
	}, nil
}

func (r *repository) BulkInsert(batch []*domain.Stock) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
    INSERT INTO stocks (ticker, trade_date, trade_price, volume, closing_time)
    VALUES %s
    ON CONFLICT (ticker, trade_date, trade_price, volume, closing_time)
    DO UPDATE SET ticker = EXCLUDED.ticker, trade_date = EXCLUDED.trade_date, trade_price = EXCLUDED.trade_price, volume = EXCLUDED.volume, closing_time = EXCLUDED.closing_time
  `
	values := []any{}
	placeholders := []string{}
	for i, stock := range batch {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		values = append(values, stock.Ticker, stock.TradeDate, stock.TradePrice, stock.Volume, stock.ClosingTime)
	}
	finalQuery := fmt.Sprintf(query, strings.Join(placeholders, ","))
	stmt, err := tx.Prepare(finalQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Query(values...)
	if err != nil {
		fmt.Printf("error inserting stock %v \n", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}

func (r *repository) Search(queryParams *schemas.StockSchema) (*domain.StockSummary, error) {
	query := `
    select ticker, max(total_volume) as max_daily_volume, max(max_price) as max_range_value from public.stock_summary
    where ticker = $1
    `
	args := []any{queryParams.Ticker}
	if queryParams.StarDate != nil {
		query = query + " and trade_date >= $2"
		args = append(args, queryParams.StarDate)
	}
	query = query + " group by ticker;"
	fmt.Println(query)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var stock domain.StockSummary
	row := stmt.QueryRow(args...)
	err = row.Scan(&stock.Ticker, &stock.MaxDailyVolume, &stock.MaxRangeValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stock info not found: %w", err)
		} else {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
	}
	return &stock, nil
}

func (r *repository) RefreshStockSummaryMeterializedView() error {
	query := `REFRESH MATERIALIZED view public.stock_summary;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(context.Background())
	if err != nil {
		return fmt.Errorf("failed to refresh materialied view: %w", err)
	}
	return nil
}
