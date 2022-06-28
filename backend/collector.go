package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rootwarp/wenmerge/rpc"
	"github.com/rs/zerolog/log"
)

var (
	wssURL = os.Getenv("ETH_WSS_URL")
)

// TODO: header buffer size?
type headCollector struct {
	ethHeaders      []*types.Header
	totalDifficulty map[string]*big.Int // Use string as key to handle big int.

	termChan chan bool
}

func (h *headCollector) Start(headRecvNotifChan chan<- *types.Header) {
	log.Info().Str("main", "headCollector.Start").Msg("start")

	h.ethHeaders = make([]*types.Header, 0)
	h.totalDifficulty = map[string]*big.Int{}

	h.termChan = make(chan bool, 1)

	go h.headReceiver(h.termChan, headRecvNotifChan)
}

func (h *headCollector) headReceiver(termChan <-chan bool, headRecvNotifChan chan<- *types.Header) error {
	log.Info().Str("main", "headCollector.headReceiver").Msg("start receiver")

	cli, err := ethclient.Dial(wssURL)
	if err != nil {
		log.Error().Str("main", "headReceiver").Msg(err.Error())
		return err
	}

	headChan := make(chan *types.Header, 2)
	ctx := context.Background()
	sub, err := cli.SubscribeNewHead(ctx, headChan)
	if err != nil {
		log.Error().Str("main", "headReceiver").Msg(err.Error())
		return err
	}

	defer sub.Unsubscribe()

	rpcCli := rpc.NewClient()

	for {
		select {
		case ethHeader := <-headChan:
			log.Info().Str("main", "headReceiver").Msg(fmt.Sprintf("%v, %v", ethHeader.Number, ethHeader.Difficulty))
			h.ethHeaders = append(h.ethHeaders, ethHeader)

			if len(h.ethHeaders) == 1 {
				block, err := rpcCli.GetBlockByNumber(ethHeader.Number)
				if err != nil {
					log.Warn().Str("main", "headReceiver").Msg(err.Error())
					continue
				}

				h.totalDifficulty[ethHeader.Number.Text(10)] = block.TotalDifficulty
			} else {
				prevBlockNo := new(big.Int).Sub(ethHeader.Number, big.NewInt(1))
				prevTotalDifficulty := h.totalDifficulty[prevBlockNo.Text(10)]
				curTotalDifficulty := new(big.Int).Add(prevTotalDifficulty, ethHeader.Difficulty)

				h.totalDifficulty[ethHeader.Number.Text(10)] = curTotalDifficulty
			}

			headRecvNotifChan <- ethHeader

			logMsg := fmt.Sprintf("%d | %d | %d", ethHeader.Number, ethHeader.Difficulty, h.totalDifficulty[ethHeader.Number.Text(10)])
			log.Info().Str("main", "headReceiver").Msg(fmt.Sprintf(logMsg))
		case _ = <-termChan:
			log.Warn().Str("main", "headReceiver").Msg("terminate")
			break
		}
	}
}

func (h *headCollector) Stop() {
	log.Info().Str("main", "headCollector.Start").Msg("stop")
	h.termChan <- true
}

func (h *headCollector) GetEthHeaders() []*types.Header {
	return h.ethHeaders
}

func (h *headCollector) GetLastEthHeader() *types.Header {
	if len(h.ethHeaders) == 0 {
		return nil
	}

	return h.ethHeaders[len(h.ethHeaders)-1]
}

func (h *headCollector) GetTotalDifficulty(blockNo *big.Int) (*big.Int, bool) {
	totalDifficulty, ok := h.totalDifficulty[blockNo.Text(10)]
	return totalDifficulty, ok
}

func (h *headCollector) GetLastTotalDifficulty() *big.Int {
	lastHeader := h.GetLastEthHeader()
	if lastHeader == nil {
		return nil
	}

	difficulty, ok := h.totalDifficulty[lastHeader.Number.Text(10)]
	if !ok {
		return nil
	}

	return difficulty
}
