package db

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// BlockHeaderStore provides interfaces both read and write from persistent service.
type BlockHeaderStore interface {
	Store(ctx context.Context, blockNo *big.Int, header *types.Header) error
	Get(ctx context.Context, blockNo *big.Int) (*types.Header, error)

	SetLatestHeader(ctx context.Context, blockNo *big.Int, header *types.Header) error
	Latest(ctx context.Context) (*types.Header, error)

	SetTotalDifficulty(ctx context.Context, blockNo *big.Int, totalDifficulty *big.Int) error
	GetTotalDifficulty(ctx context.Context, blockNo *big.Int) (*big.Int, error)
}

// NewClient creates persistent layer client.
func NewClient(addr string) BlockHeaderStore {
	return newRedisClient(addr)
}
