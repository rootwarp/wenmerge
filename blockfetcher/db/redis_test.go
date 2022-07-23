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
// TODO: Check empty value.

// TODO: Sort by big.int?

// Do the simple thing.
// Get all and sort.

//
// NOTICE: For now, you need to run Redis server on local machine.
//
func TestRedisSet(t *testing.T) {
	store := newRedisClient("localhost:6379")

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

func TestTotalDifficulty(t *testing.T) {
	store := newRedisClient("localhost:6379")

	ctx := context.Background()
	err := store.SetTotalDifficulty(ctx, big.NewInt(1234), big.NewInt(4321))

	assert.Nil(t, err)

	readDifficulty, err := store.GetTotalDifficulty(ctx, big.NewInt(1234))

	assert.Nil(t, err)
	assert.Equal(t, "4321", readDifficulty.String())
}

func TestLatest(t *testing.T) {
	store := newRedisClient("localhost:6379")

	// dummy info.
	header := &types.Header{
		Coinbase:   common.HexToAddress("12345"),
		Number:     big.NewInt(1234),
		Difficulty: big.NewInt(56789),
		Time:       uint64(time.Now().Unix()),
	}

	// Set
	ctx := context.Background()
	err := store.SetLatestHeader(ctx, header.Number, header)

	assert.Nil(t, err)

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// Check latest block number value.
	latestBlockNo, err := rdb.Get(ctx, keyLatestBlockNo).Result()

	assert.Nil(t, err)
	assert.Equal(t, header.Number.String(), latestBlockNo)

	// Read real header.
	recvHeader, err := store.Get(ctx, header.Number)

	assert.Nil(t, err)
	assert.Equal(t, header.Number.String(), recvHeader.Number.String())
}

func TestEmpty(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	ctx := context.Background()
	cmd := rdb.Get(ctx, "no_key")

	value, err := cmd.Result()

	assert.Empty(t, value)
	assert.NotNil(t, err)
}

func TestGetLatest(t *testing.T) {
	// Set dummies.
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	store := newRedisClient("localhost:6379")

	// Store
	header := &types.Header{
		Coinbase:   common.HexToAddress("12345"),
		Number:     big.NewInt(1234),
		Difficulty: big.NewInt(56789),
		Time:       uint64(time.Now().Unix()),
	}

	_ = store.Store(ctx, header.Number, header)
	rdb.Set(ctx, keyLatestBlockNo, header.Number.String(), 0)

	// Get
	readHeader, err := store.Latest(ctx)

	// Asserts
	assert.Nil(t, err)
	assert.Equal(t, header.Number.String(), readHeader.Number.String())
	assert.Equal(t, header.Difficulty.String(), readHeader.Difficulty.String())
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
	assert.Nil(t, err)

	assert.Equal(t, header.Number, newInfo.Number)
	assert.Equal(t, header.Difficulty, newInfo.Difficulty)
	assert.Equal(t, header.Time, newInfo.Time)
}
