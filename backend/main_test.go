package main

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
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

func TestFetchAhead(t *testing.T) {
	collector := headCollector{
		ethHeaders:      make([]*types.Header, 0),
		totalDifficulty: map[string]*big.Int{},
		headerLock:      sync.Mutex{},
		difficultyLock:  sync.Mutex{},
		termChan:        make(chan bool, 1),
	}

	blockNo := big.NewInt(1000000)
	collector.fetchAheadHeaders(blockNo, 10)
}
