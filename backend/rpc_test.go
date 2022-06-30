package main

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	rpcURL = "https://ethereum-mainnet-rpc.allthatnode.com/qUH1yZIx5qxXSj6Pi7ZHLLOlKnZMcE6i"
)

func TestGetBlockNumber(t *testing.T) {
	cli, err := ethclient.Dial(rpcURL)

	ctx := context.Background()
	no, err := cli.BlockNumber(ctx)

	fmt.Println(no, err)

}

func TestGetBlockByNumber(t *testing.T) {
	infuraRPC := "https://mainnet.infura.io/v3/b9c3f7e7460345bdb661add3b5e451c7"
	cli, err := ethclient.Dial(infuraRPC)

	ctx := context.Background()
	block, err := cli.BlockByNumber(ctx, big.NewInt(100000))

	fmt.Println(block.Number(), block.Header().Hash())

	_ = err
}
