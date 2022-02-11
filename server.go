package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port int `required:"true"`
	Pg   struct {
		Network      string
		Addr         string
		Database     string `required:"true"`
		TestDatabase string `split_words:"true"`
		User         string `required:"true"`
		Password     string `required:"true"`
	}
}

type Handler struct {
	symbolRepo stockSymbolRepository
	priceRepo  stockPriceRepository
}

var handler Handler

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	pgOpts := &pg.Options{
		Network:  cfg.Pg.Network,
		Addr:     cfg.Pg.Addr,
		Database: cfg.Pg.Database,
		User:     cfg.Pg.User,
		Password: cfg.Pg.Password,
	}
	db := pg.Connect(pgOpts)
	defer db.Close()

	handler = Handler{&pgStockSymbolRepository{db}, &pgStockPriceRepository{db}}

	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Logger())
	e.GET("/", handler.getHome)
	e.GET("symbols", handler.getSymbols)
	e.GET("/prices/:symbol", handler.getPricesBySymbol)
	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(cfg.Port)))
}

func (h Handler) getHome(ctx echo.Context) error {
	tmpl, err := template.ParseGlob("views/*.html")
	if err != nil {
		return err
	}

	sb := &strings.Builder{}
	tmpl.ExecuteTemplate(sb, "chart", nil)
	return ctx.HTML(http.StatusOK, sb.String())
}

func (h Handler) getSymbols(ctx echo.Context) error {
	symbols, err := h.symbolRepo.get()
	if err != nil {
		return err
	}

	syms := make([]string, len(symbols))
	for i, symbol := range symbols {
		syms[i] = symbol.Symbol
	}
	return ctx.JSON(http.StatusOK, wrap(syms))
}

func (h Handler) getPricesBySymbol(ctx echo.Context) error {
	symbol := strings.ToUpper(ctx.Param("symbol")) + ".XIDX"
	prices, err := h.priceRepo.getBySymbol(symbol)
	if err != nil {
		return err
	}

	if prices == nil {
		return ctx.NoContent(http.StatusNotFound)
	}

	return ctx.JSON(http.StatusOK, wrap(prices))
}

func wrap(v interface{}) map[string]interface{} {
	return map[string]interface{}{"data": v}
}
