package db

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
)

// TODO: Get latest.
// TODO: Sort
// TODO: Mock test.

// TODO: Sort by big.int?

// Do the simple thing.
// Get all and sort.

//
// NOTICE: For now, you need to run Redis server on local machine.
//
func TestRedisSet(t *testing.T) {
	store := newRedisClient("localhost:6379")

	// Connect
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	// Store
	header := &types.Header{
		Coinbase:   common.HexToAddress("12345"),
		Number:     big.NewInt(1234),
		Difficulty: big.NewInt(56789),
		Time:       uint64(time.Now().Unix()),
	}

	ctx := context.Background()

	err := store.Store(ctx, header.Number, header)
	assert.Nil(t, err)

	getHeader, err := store.Get(ctx, header.Number)
	assert.Nil(t, err)

	assert.Equal(t, header.Number.String(), getHeader.Number.String())
	assert.Equal(t, header.Coinbase.String(), getHeader.Coinbase.String())
	assert.Equal(t, header.Difficulty.String(), getHeader.Difficulty.String())
	assert.Equal(t, header.Time, getHeader.Time)
}

func TestGob(t *testing.T) {
	store := redisStore{}

	header := &types.Header{
		Coinbase:   common.HexToAddress("12345"),
		Number:     big.NewInt(1234),
		Difficulty: big.NewInt(56789),
		Time:       uint64(time.Now().Unix()),
	}

	hexStr, err := store.encode(header)
	assert.Nil(t, err)

	newInfo, err := store.decode(hexStr)

	assert.Equal(t, header.Number, newInfo.Number)
	assert.Equal(t, header.Difficulty, newInfo.Difficulty)
	assert.Equal(t, header.Time, newInfo.Time)
}
