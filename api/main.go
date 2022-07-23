package main

import (
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	reader = NewBlockReader()
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

func main() {
	e := echo.New()

	corsConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodHead},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}

	e.Use(middleware.CORSWithConfig(corsConfig))

	e.GET("/difficulty", diffHandler)

	if err := e.Start("0.0.0.0:9090"); err != nil {
		panic(err)
	}
}

func diffHandler(c echo.Context) error {
	log.Info().Str("module", "main").Msg("diffHandler")

	if c.QueryParam("target") == "" {
		return c.String(http.StatusBadRequest, "need target parameter")
	}

	targetDifficulty, ok := new(big.Int).SetString(c.QueryParam("target"), 10)
	if !ok {
		return c.String(http.StatusBadRequest, "invalid target difficulty number format")
	}

	ctx := c.Request().Context()
	latestHeader, err := reader.Latest(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	aheadHeaderNo := new(big.Int).Sub(latestHeader.Number, big.NewInt(500))

	aheadHeader, err := reader.Get(ctx, aheadHeaderNo)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	difficultyDiff := new(big.Int).Sub(latestHeader.TotalDifficulty, aheadHeader.TotalDifficulty)
	diffBlockNo := new(big.Int).Sub(latestHeader.Number, aheadHeader.Number)
	difficultyVelocity := new(big.Int).Div(difficultyDiff, diffBlockNo)

	ttdDistance := new(big.Int).Sub(targetDifficulty, latestHeader.TotalDifficulty)
	remainBlocks := new(big.Int).Div(ttdDistance, difficultyVelocity)

	avgBlockSec := float64(latestHeader.Time-aheadHeader.Time) / float64(diffBlockNo.Int64())
	expectTTDTime := time.Now().Add(time.Millisecond * time.Duration(avgBlockSec*1000))

	respBody := struct {
		CurrentBlockNumber    *big.Int  `json:"current_block_number"`
		CurrentDifficulty     *big.Int  `json:"current_difficulty"`
		CurrentBlockTimestamp time.Time `json:"current_block_timestamp"`
		TargetDifficulty      *big.Int  `json:"target_difficulty"`
		DifficultyVelocity    *big.Int  `json:"difficulty_velocity"`
		ExpectTTDBlockNumber  *big.Int  `json:"expect_ttd_block_number"`
		ExpectTTDTime         time.Time `json:"expect_ttd_time"`
		AverageBlockInterval  float64   `json:"average_block_interval"`
	}{
		CurrentBlockNumber:    latestHeader.Number,
		CurrentDifficulty:     latestHeader.Difficulty,
		CurrentBlockTimestamp: time.Unix(int64(latestHeader.Time), 0),
		TargetDifficulty:      targetDifficulty,
		DifficultyVelocity:    difficultyVelocity,
		ExpectTTDBlockNumber:  new(big.Int).Add(latestHeader.Number, remainBlocks),
		ExpectTTDTime:         expectTTDTime,
		AverageBlockInterval:  avgBlockSec,
	}

	return c.JSON(http.StatusOK, respBody)
}
