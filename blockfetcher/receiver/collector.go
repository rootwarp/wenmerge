package receiver

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"

	"github.com/rootwarp/wenmerge/blockfetcher/db"
	"github.com/rootwarp/wenmerge/blockfetcher/rpc"
)

var (
	wssURL = os.Getenv("ETH_WSS_URL")
	rpcURL = os.Getenv("ETH_RPC_URL")
)

// Headers is type of header array.
type Headers []*types.Header

func (h Headers) Len() int           { return len(h) }
func (h Headers) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Headers) Less(i, j int) bool { return h[i].Number.Cmp(h[j].Number) < 0 }

// HeaderReceiver controls Ether header receiver.
type HeaderReceiver interface {
	Start(headRecvNotifChan chan<- *types.Header)
	Stop()
}

type headerReceiver struct {
	termChan chan bool

	store db.BlockHeaderStore
}

// TODO: need notif chan?
func (h *headerReceiver) Start(headRecvNotifChan chan<- *types.Header) {
	log.Info().Str("module", "receiver").Msg("start")

	h.termChan = make(chan bool, 1)

	// Get Latest
	//go h.fetchAheadHeaders(ethHeader.Number, 10) // TODO

	ticker := time.NewTicker(time.Second * 60)

	go h.headReceiver(h.termChan)
	for t := range ticker.C {
		h.termChan <- true
		time.Sleep(time.Second * 1)

		go h.headReceiver(h.termChan)
		_ = t
	}
}

func (h *headerReceiver) headReceiver(termChan <-chan bool) error {
	log.Info().Str("module", "receiver").Msg("start headReceiver")

	cli, err := ethclient.Dial(wssURL)
	if err != nil {
		log.Error().Str("module", "receiver").Msg(err.Error())
		return err
	}

	headChan := make(chan *types.Header, 100)
	ctx := context.Background()
	sub, err := cli.SubscribeNewHead(ctx, headChan)
	if err != nil {
		log.Error().Str("module", "receiver").Msg(err.Error())
		return err
	}

	defer sub.Unsubscribe()

	rpcCli := rpc.NewClient()

	for {
		select {
		case ethHeader := <-headChan:
			log.Info().Str("module", "receiver").Msg(fmt.Sprintf("got %v, %v", ethHeader.Number, ethHeader.Difficulty))

			err = h.store.SetLatestHeader(ctx, ethHeader.Number, ethHeader)
			if err != nil {
				log.Error().Str("module", "receiver").Msg(err.Error())
				continue
			}

			block, err := rpcCli.GetBlockByNumber(ethHeader.Number)
			if err != nil {
				log.Error().Str("module", "receiver").Msg(err.Error())
				continue
			}

			err = h.store.SetTotalDifficulty(ctx, ethHeader.Number, block.TotalDifficulty)
			if err != nil {
				log.Error().Str("module", "receiver").Msg(err.Error())
				continue
			}
		case err := <-sub.Err():
			log.Error().Str("module", "receiver").Msg(err.Error())

			cli, err = ethclient.Dial(wssURL)
			if err != nil {
				log.Error().Str("module", "receiver").Msg(err.Error())
				return err
			}

			sub, err = cli.SubscribeNewHead(ctx, headChan)
			if err != nil {
				log.Error().Str("module", "receiver").Msg(err.Error())
				return err
			}
		case _ = <-termChan:
			log.Warn().Str("module", "receiver").Msg("terminate")
			return nil
		}
	}
}

func (h *headerReceiver) fetchAheadHeaders(blockNo *big.Int, count int) {
	log.Info().Str("main", "headerReceiver.fetchAheadHeaders").Msgf("fetch %d from %d", count, blockNo)

	reqBlockChan := make(chan *big.Int, 100)

	newHeaders := make([]*types.Header, 0)
	lock := sync.Mutex{}

	const nWorkers = 5
	wg := sync.WaitGroup{}
	wg.Add(nWorkers)

	for i := 0; i < nWorkers; i++ {
		go func(c <-chan *big.Int) {
			defer wg.Done()
			cli, err := ethclient.Dial(rpcURL)
			if err != nil {
				log.Error().Str("main", "headerReceiver.fetchAheadHeaders").Msg(err.Error())
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
					log.Error().Str("main", "headerReceiver.fetchAheadHeaders").Msg(err.Error())
					return
				}

				log.Info().Str("main", "headerReceiver.fetchAheadHeaders").Msgf("Append %d", block.Number())
				lock.Lock()
				newHeaders = append(newHeaders, block.Header())
				lock.Unlock()

				err = h.store.Store(ctx, block.Number(), block.Header())
				if err != nil {
					log.Error().Str("main", "headerReceiver.fetchAheadHeaders").Msg(err.Error())
					return
				}
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
}

func (h *headerReceiver) Stop() {
	log.Info().Str("module", "receiver").Msg("stop")
	h.termChan <- true
}

// NewReceiver creates new receiver instance.
func NewReceiver() HeaderReceiver {
	redisAddr := os.Getenv("REDIS_ADDR")
	return &headerReceiver{
		store: db.NewClient(redisAddr),
	}
}
