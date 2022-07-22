package main

import (
	"github.com/ethereum/go-ethereum/core/types"

	recv "github.com/rootwarp/wenmerge/blockfetcher/receiver"
)

func main() {
	// TODO: Start fetcher.

	//
	// TODO: Preload.
	// Get current latest.
	// Get current Redis latest.
	// Get blocks.
	//

	headerChan := make(chan *types.Header, 10)

	receiver := recv.NewReceiver()
	receiver.Start(headerChan)
}

func getBlocks() []*types.Header {
	return nil
}
