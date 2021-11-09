package main

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type stockPrice struct {
	tableName struct{}  `pg:"daily_stock_prices"`
	Open      float32   `json:"open"`
	Low       float32   `json:"low"`
	High      float32   `json:"high"`
	Close     float32   `json:"close"`
	Time      time.Time `json:"time" pg:"date"`
}

type stockPriceRepository interface {
	getBySymbol(string) ([]stockPrice, error)
}

type pgStockPriceRepository struct {
	db *pg.DB
}

func (repo *pgStockPriceRepository) getBySymbol(symbol string) (prices []stockPrice, err error) {
	err = repo.db.Model(&prices).Where("symbol = ?", symbol).Order("date").Select()
	return prices, err
}
