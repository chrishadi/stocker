package main

import (
	"github.com/go-pg/pg/v10"
)

type stockSymbol struct {
	tableName struct{} `pg:"stock_last_updates"`
	Symbol    string   `pg:"code"`
}

type stockSymbolRepository interface {
	get() ([]stockSymbol, error)
}

type pgStockSymbolRepository struct {
	db *pg.DB
}

func (repo *pgStockSymbolRepository) get() (symbols []stockSymbol, err error) {
	err = repo.db.Model(&symbols).Order("code").Select()
	return symbols, err
}
