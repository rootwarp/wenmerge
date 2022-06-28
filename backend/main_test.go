package main

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDiffHandler(t *testing.T) {
	// Prepare
	e := echo.New()

	q := make(url.Values)
	q.Set("target", "52836506573602732113920")

	req := httptest.NewRequest(http.MethodGet, "/Difficulty?"+q.Encode(), nil)
	rr := httptest.NewRecorder()

	c := e.NewContext(req, rr)

	// Test
	err := diffHandler(c)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, rr.Result().StatusCode)
}

func TestSubBigInt(t *testing.T) {
	z := new(big.Int)

	x := big.NewInt(10000)
	y := big.NewInt(1)

	z = z.Sub(x, y)

	fmt.Println(z)
}
