package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	keyLatestBlockNo = "latest_block_no"
)

// BlockHeader privides wrap properties of original Eth header.
type BlockHeader struct {
	Number          *big.Int `json:"number"`
	Difficulty      *big.Int `json:"difficulty"`
	TotalDifficulty *big.Int `json:"total_difficulty"`
	Time            uint64   `json:"time"`
}

// BlockHeaderReader provides reader interfaces to read Eth block.
type BlockHeaderReader interface {
	Latest(ctx context.Context) (*BlockHeader, error)
	Get(ctx context.Context, blockNo *big.Int) (*BlockHeader, error)
}

type redisReader struct {
	redisCli *redis.Client
}

func (r *redisReader) Latest(ctx context.Context) (*BlockHeader, error) {
	log.Info().Str("module", "redis").Msg("latest")

	blockNoStr, err := r.redisCli.Get(ctx, keyLatestBlockNo).Result()
	if err != nil {
		log.Error().Str("module", "redis").Msgf("err get latest block no. %+v", err)
		return nil, err
	}

	blockNo, ok := big.NewInt(0).SetString(blockNoStr, 10)
	if !ok {
		return nil, errors.New("cannot convert string to big.Int")
	}

	return r.Get(ctx, blockNo)
}

func (r *redisReader) decode(encoded string) (*types.Header, error) {
	rawData, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewBuffer(rawData)
	dec := gob.NewDecoder(reader)

	info := types.Header{}
	err = dec.Decode(&info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (r *redisReader) Get(ctx context.Context, blockNo *big.Int) (*BlockHeader, error) {
	log.Info().Str("module", "redis").Msgf("get %s", blockNo.String())

	encBlock, err := r.redisCli.Get(ctx, blockNo.String()).Result()
	if err != nil {
		log.Error().Str("module", "redis").Msgf("err get block: %+v", err)
		return nil, err
	}

	header, err := r.decode(encBlock)
	if err != nil {
		log.Error().Str("module", "redis").Msgf("err decode: %+v", err)
		return nil, err
	}

	totalDifficultyStr, err := r.redisCli.Get(ctx, "td_"+blockNo.String()).Result()
	if err != nil {
		log.Error().Str("module", "redis").Msgf("err get total difficulty: %+v", err)
		return nil, err
	}

	totalDifficulty, ok := big.NewInt(0).SetString(totalDifficultyStr, 10)
	if !ok {
		return nil, errors.New("cannot convert string to big.Int")
	}

	blockHeader := BlockHeader{
		Number:          header.Number,
		Difficulty:      header.Difficulty,
		TotalDifficulty: totalDifficulty,
		Time:            header.Time,
	}

	return &blockHeader, nil
}

// NewBlockReader creates reader instance for reading eth block header.
func NewBlockReader() BlockHeaderReader {
	cli := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR")})
	return &redisReader{redisCli: cli}
}
