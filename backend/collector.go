package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rootwarp/wenmerge/rpc"
	"github.com/rs/zerolog/log"
)

var (
	wssURL = os.Getenv("ETH_WSS_URL")
)

// Headers is type of header array.
type Headers []*types.Header

func (h Headers) Len() int           { return len(h) }
func (h Headers) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Headers) Less(i, j int) bool { return h[i].Number.Cmp(h[j].Number) < 0 }

// TODO: header buffer size?
type headCollector struct {
	ethHeaders      []*types.Header
	totalDifficulty map[string]*big.Int // Use string as key to handle big int.

	headerLock     sync.Mutex
	difficultyLock sync.Mutex

	termChan chan bool
}

func (h *headCollector) Start(headRecvNotifChan chan<- *types.Header) {
	log.Info().Str("main", "headCollector.Start").Msg("start")

	h.ethHeaders = make([]*types.Header, 0)
	h.totalDifficulty = map[string]*big.Int{}

	h.headerLock = sync.Mutex{}
	h.difficultyLock = sync.Mutex{}

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

			h.headerLock.Lock()
			h.ethHeaders = append(h.ethHeaders, ethHeader)
			h.headerLock.Unlock()

			if len(h.ethHeaders) == 1 {
				block, err := rpcCli.GetBlockByNumber(ethHeader.Number)
				if err != nil {
					log.Warn().Str("main", "headReceiver").Msg(err.Error())
					continue
				}

				h.difficultyLock.Lock()
				h.totalDifficulty[ethHeader.Number.Text(10)] = block.TotalDifficulty
				h.difficultyLock.Unlock()

				go h.fetchAheadHeaders(ethHeader.Number, 10)
			} else {
				prevBlockNo := new(big.Int).Sub(ethHeader.Number, big.NewInt(1))
				prevTotalDifficulty := h.totalDifficulty[prevBlockNo.Text(10)]
				curTotalDifficulty := new(big.Int).Add(prevTotalDifficulty, ethHeader.Difficulty)

				h.difficultyLock.Lock()
				h.totalDifficulty[ethHeader.Number.Text(10)] = curTotalDifficulty
				h.difficultyLock.Unlock()
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

func (h *headCollector) fetchAheadHeaders(blockNo *big.Int, count int) {
	log.Info().Str("main", "headCollector.fetchAheadHeaders").Msgf("fetch %d from %d", count, blockNo)

	reqBlockChan := make(chan *big.Int, 100)

	newHeaders := make([]*types.Header, 0)
	lock := sync.Mutex{}

	const nWorkers = 5
	wg := sync.WaitGroup{}
	wg.Add(nWorkers)

	for i := 0; i < nWorkers; i++ {
		go func(c <-chan *big.Int) {
			defer wg.Done()
			cli, err := ethclient.Dial(wssURL)
			if err != nil {
				log.Error().Str("main", "headCollector.fetchAheadHeaders").Msg(err.Error())
				return
			}

			defer cli.Close()

			ctx := context.Background()

			for {
				n, ok := <-c
				if !ok {
					break
				}

				block, err := cli.BlockByNumber(ctx, n)
				if err != nil {
					log.Error().Str("main", "headCollector.fetchAheadHeaders").Msg(err.Error())
					return
				}

				log.Info().Str("main", "headCollector.fetchAheadHeaders").Msgf("Append %d", block.Number())
				lock.Lock()
				newHeaders = append(newHeaders, block.Header())
				lock.Unlock()
			}

		}(reqBlockChan)
	}

	for i := 0; i < count; i++ {
		number := new(big.Int).Sub(blockNo, big.NewInt(int64(i+1)))
		reqBlockChan <- number
	}

	close(reqBlockChan)
	wg.Wait()

	sort.Sort(Headers(newHeaders))

	// Update total difficult.
	h.headerLock.Lock()
	defer h.headerLock.Unlock()

	h.difficultyLock.Lock()
	defer h.difficultyLock.Unlock()

	rpcCli := rpc.NewClient()
	block, err := rpcCli.GetBlockByNumber(newHeaders[0].Number)
	if err != nil {
		log.Warn().Str("main", "headReceiver").Msg(err.Error())
	}

	h.totalDifficulty[block.Number.Text(10)] = block.TotalDifficulty
	for _, b := range newHeaders[1:] {
		prevBlockNo := new(big.Int).Sub(b.Number, big.NewInt(1))
		prevTotalDifficulty, ok := h.totalDifficulty[prevBlockNo.Text(10)]
		if !ok {
			// TODO:
			panic("can't find...")
		}

		h.totalDifficulty[b.Number.Text(10)] = new(big.Int).Add(prevTotalDifficulty, b.Difficulty)
		fmt.Println(b.Number, b.Difficulty, new(big.Int).Add(prevTotalDifficulty, b.Difficulty))
	}

	h.ethHeaders = append(newHeaders, h.ethHeaders...)

	for _, head := range newHeaders {
		fmt.Println(head.Number, h.totalDifficulty[head.Number.Text(10)])
	}
}

func (h *headCollector) Stop() {
	log.Info().Str("main", "headCollector.Start").Msg("stop")
	h.termChan <- true
}

func (h *headCollector) GetEthHeaders() []*types.Header {
	h.headerLock.Lock()
	defer h.headerLock.Unlock()

	return h.ethHeaders
}

func (h *headCollector) GetLastEthHeader() *types.Header {
	h.headerLock.Lock()
	defer h.headerLock.Unlock()

	if len(h.ethHeaders) == 0 {
		return nil
	}

	return h.ethHeaders[len(h.ethHeaders)-1]
}

func (h *headCollector) GetTotalDifficulty(blockNo *big.Int) (*big.Int, bool) {
	h.difficultyLock.Lock()
	defer h.difficultyLock.Unlock()

	totalDifficulty, ok := h.totalDifficulty[blockNo.Text(10)]
	return totalDifficulty, ok
}

func (h *headCollector) GetLastTotalDifficulty() *big.Int {
	h.difficultyLock.Lock()
	defer h.difficultyLock.Unlock()

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
