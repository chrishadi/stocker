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

	e := echo.New()
	e.Use(middleware.Gzip())
	e.GET("/", func(ctx echo.Context) error {
		tmpl := template.Must(template.ParseGlob("views/*.html"))
		sb := &strings.Builder{}
		tmpl.ExecuteTemplate(sb, "chart", nil)
		return ctx.HTML(http.StatusOK, sb.String())
	})
	e.GET("symbols", func(ctx echo.Context) error {
		db := pg.Connect(pgOpts)
		repo := &pgStockSymbolRepository{db: db}
		symbols, err := repo.get()
		if err != nil {
			e.Logger.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		syms := make([]string, len(symbols))
		for i, symbol := range symbols {
			syms[i] = symbol.Symbol
		}
		return ctx.JSON(http.StatusOK, syms)
	})
	e.GET("/prices/:symbol", func(ctx echo.Context) error {
		symbol := strings.ToUpper(ctx.Param("symbol")) + ".XIDX"
		db := pg.Connect(pgOpts)
		repo := &pgStockPriceRepository{db: db}
		prices, err := repo.getBySymbol(symbol)
		if err != nil {
			e.Logger.Error(err)
			return ctx.NoContent(http.StatusInternalServerError)
		}
		if prices == nil {
			return ctx.NoContent(http.StatusNotFound)
		}
		return ctx.JSON(http.StatusOK, prices)
	})
	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(cfg.Port)))
}
