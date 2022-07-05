package main

import (
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	minBlocksToEstimate = 10
)

var (
	collector          headCollector
	velocityCalculator difficultyVelocityCalculator
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

func main() {
	headRecvNotifChan := make(chan *types.Header, 10)

	// TODO: Should be singleton.
	collector = headCollector{}
	collector.Start(headRecvNotifChan)

	velocityCalculator = difficultyVelocityCalculator{}
	velocityCalculator.Start(headRecvNotifChan)

	e := echo.New()
	e.GET("/difficulty", diffHandler)
	e.GET("/stat", statHandler)

	corsConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodHead},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}

	e.Use(middleware.CORSWithConfig(corsConfig))

	if err := e.Start("0.0.0.0:9090"); err != nil {
		collector.Stop()
		panic(err)
	}
}

func diffHandler(c echo.Context) error {
	log.Info().Str("main", "diffHandler").Msg("call")

	if c.QueryParam("target") == "" {
		return c.String(http.StatusBadRequest, "need target parameter")
	}

	if len(collector.GetEthHeaders()) < minBlocksToEstimate {
		return c.String(http.StatusNoContent, "not ready")
	}

	lastEthHeader := collector.GetLastEthHeader()
	if lastEthHeader == nil {
		return c.String(http.StatusInternalServerError, "can't get last header")
	}

	targetDifficulty, ok := new(big.Int).SetString(c.QueryParam("target"), 10)
	if !ok {
		return c.String(http.StatusBadRequest, "invalid target difficulty number format")
	}

	// TODO: Handle past..

	curDifficulty := collector.GetLastTotalDifficulty()
	if curDifficulty == nil {
		return c.String(http.StatusInternalServerError, "can't get last difficulty")
	}

	distance := new(big.Int).Sub(targetDifficulty, curDifficulty)
	expectBlocks := new(big.Int).Div(distance, velocityCalculator.velocity)

	firstEthHeader := collector.GetEthHeaders()[0]
	diffBlockNo := new(big.Int).Sub(lastEthHeader.Number, firstEthHeader.Number).Uint64()
	avgBlockInterval := (lastEthHeader.Time - firstEthHeader.Time) / diffBlockNo

	var ethBlockInterval = time.Duration(avgBlockInterval) * time.Second
	now := time.Now().UTC()
	expectTTDTime := now.Add(ethBlockInterval * time.Duration(expectBlocks.Int64()))

	fmt.Println(diffBlockNo, (lastEthHeader.Time - firstEthHeader.Time), avgBlockInterval)

	respBody := struct {
		CurrentBlockNumber    *big.Int  `json:"current_block_number"`
		CurrentDifficulty     *big.Int  `json:"current_difficulty"`
		CurrentBlockTimestamp time.Time `json:"current_block_timestamp"`
		TargetDifficulty      *big.Int  `json:"target_difficulty"`
		DifficultyVelocity    *big.Int  `json:"difficulty_velocity"`
		ExpectTTDBlockNumber  *big.Int  `json:"expect_ttd_block_number"`
		ExpectTTDTime         time.Time `json:"expect_ttd_time"`
		AverageBlockInterval  uint64    `json:"average_block_interval"`
	}{
		CurrentBlockNumber:    lastEthHeader.Number,
		CurrentDifficulty:     lastEthHeader.Difficulty,
		CurrentBlockTimestamp: time.Unix(int64(lastEthHeader.Time), 0),
		TargetDifficulty:      targetDifficulty,
		DifficultyVelocity:    velocityCalculator.velocity, // TODO: wrap
		ExpectTTDBlockNumber:  new(big.Int).Add(lastEthHeader.Number, expectBlocks),
		ExpectTTDTime:         expectTTDTime,
		AverageBlockInterval:  avgBlockInterval,
	}

	return c.JSON(http.StatusOK, respBody)
}

func statHandler(c echo.Context) error {
	fmt.Println("Statistics")
	fmt.Printf("Total %d headers\n", len(collector.GetEthHeaders()))
	fmt.Printf("Total %d diff\n", len(collector.totalDifficulty))

	// TODO: Handle by response.
	for _, h := range collector.GetEthHeaders() {
		difficulty := h.Difficulty
		totalDifficulty := collector.totalDifficulty[h.Number.Text(10)]

		fmt.Println(h.Number, difficulty, totalDifficulty)
	}

	return c.String(http.StatusOK, "OK")
}
