package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupEndpoint(path string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestGetHomeWhenNoErrorShouldHaveHttpStatusOK(t *testing.T) {
	c, rec := setupEndpoint("/")

	if assert.NoError(t, handler.getHome(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

type MockSymbolRepository struct {
	GetResult []stockSymbol
	GetError  error
}

func (repo *MockSymbolRepository) get() ([]stockSymbol, error) {
	return repo.GetResult, repo.GetError
}

func read(r io.Reader) []byte {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return buf
}

func encode(v interface{}) []byte {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(v)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestGetSymbolsWhenNoErrorShouldRespondWithListOfSymbols(t *testing.T) {
	symbols := []stockSymbol{{Symbol: "AABB"}, {Symbol: "BBAA"}}
	mockSymbolRepo := MockSymbolRepository{symbols, nil}
	handler.symbolRepo = &mockSymbolRepo
	c, rec := setupEndpoint("/symbol")

	if assert.NoError(t, handler.getSymbols(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := read(rec.Body)
		expected := encode(map[string][]string{"data": {"AABB", "BBAA"}})
		assert.Equal(t, expected, actual)
	}
}

func TestGetSymbolsWhenErrorShouldRespondWithErrorMessage(t *testing.T) {
	mockSymbolRepo := MockSymbolRepository{nil, errors.New("symbol-error")}
	handler.symbolRepo = &mockSymbolRepo
	c, rec := setupEndpoint("/symbol")

	if assert.Error(t, handler.getSymbols(c)) {
		assert.Contains(t, "{message:", rec.Body.String())
	}
}

type MockPriceRepository struct {
	GetBySymbolResult []stockPrice
	GetBySymbolError  error
}

func (repo *MockPriceRepository) getBySymbol(symbol string) ([]stockPrice, error) {
	return repo.GetBySymbolResult, repo.GetBySymbolError
}

func TestGetPricesBySymbolWhenNoErrorShouldRespondWithPrices(t *testing.T) {
	prices := []stockPrice{
		{Open: 1.0, Low: 0.9, High: 1.2, Close: 1.1, Time: time.Now().AddDate(0, 0, -1)},
		{Open: 1.2, Low: 1.0, High: 1.4, Close: 1.4, Time: time.Now()},
	}
	mockPriceRepo := MockPriceRepository{prices, nil}
	handler.priceRepo = &mockPriceRepo
	c, rec := setupEndpoint("/prices/AABB")
	c.SetParamNames("symbol")
	c.SetParamValues("AABB")

	if assert.NoError(t, handler.getPricesBySymbol(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		actual := read(rec.Body)
		expected := encode(map[string][]stockPrice{"data": prices})
		assert.Equal(t, expected, actual)
	}
}

func TestGetPricesBySymbolWhenNoDataShouldRespondWithNotFound(t *testing.T) {
	mockPriceRepo := MockPriceRepository{nil, nil}
	handler.priceRepo = &mockPriceRepo
	c, rec := setupEndpoint("/prices/AABB")
	c.SetParamNames("symbol")
	c.SetParamValues("AABB")

	if assert.NoError(t, handler.getPricesBySymbol(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestGetPriceBySymbolWhenErrorShouldRespondWithErrorMessage(t *testing.T) {
	mockPriceRepo := MockPriceRepository{nil, errors.New("price-error")}
	handler.priceRepo = &mockPriceRepo
	c, rec := setupEndpoint("/prices/AABB")

	if assert.Error(t, handler.getPricesBySymbol(c)) {
		assert.Contains(t, "{message:", rec.Body.String())
	}
}
