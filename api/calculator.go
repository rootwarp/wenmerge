package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type difficultyVelocityCalculator struct {
	velocity      *big.Int
	blockInterval time.Duration
}

func (v *difficultyVelocityCalculator) Start(headRecvNotifChan <-chan *types.Header) {
	log.Info().Str("main", "difficultyVelocityCalculator.Start").Msg("start velocity calculator")

	go func() {
		for {
			ethHeader := <-headRecvNotifChan
			log.Info().Str("main", "difficultyVelocityCalculator").Msg(fmt.Sprintf("Get Header %v", ethHeader.Number))

			if len(collector.GetEthHeaders()) < 2 {
				log.Info().Str("main", "difficultyVelocityCalculator").Msg(fmt.Sprintf("insufficient headers"))
				continue
			}

			allHeaders := collector.GetEthHeaders()
			diffBlockNo := new(big.Int).Sub(ethHeader.Number, allHeaders[0].Number)

			lastTotalDifficulty, ok := collector.GetTotalDifficulty(ethHeader.Number)
			if !ok {
				log.Warn().Str("main", "difficultyVelocityCalculator.Start").Msg("can not find block")
				continue
			}

			firstTotalDifficulty, ok := collector.GetTotalDifficulty(collector.GetEthHeaders()[0].Number)
			if !ok {
				log.Warn().Str("main", "difficultyVelocityCalculator.Start").Msg("can not find block")
				continue
			}

			diffDifficulty := new(big.Int).Sub(lastTotalDifficulty, firstTotalDifficulty)

			velocity := new(big.Int).Div(diffDifficulty, diffBlockNo)
			v.velocity = velocity

			log.Info().Str("main", "difficultyVelocityCalculator.Start").Msgf("Velocity %d", velocity)
		}
	}()
}
