package db

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	keyLatestBlockNo         = "latest_block_no"
	keyTotalDifficultyPrefix = "td_"
)

type redisStore struct {
	client *redis.Client
}

func (s *redisStore) Store(ctx context.Context, blockNo *big.Int, header *types.Header) error {
	log.Info().Str("module", "redis").Msgf("store %s", blockNo.String())

	hexInfo, err := s.encode(header)
	if err != nil {
		return err
	}

	const expiration = time.Hour * 24 * 7

	err = s.client.Set(ctx, header.Number.String(), hexInfo, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *redisStore) encode(header *types.Header) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(header)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

func (s *redisStore) Get(ctx context.Context, blockNo *big.Int) (*types.Header, error) {
	log.Info().Str("module", "redis").Msgf("get %s", blockNo.String())

	value, err := s.client.Get(ctx, blockNo.String()).Result()
	if err != nil {
		return nil, err
	}

	return s.decode(value)
}

func (s *redisStore) decode(encoded string) (*types.Header, error) {
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

func (s *redisStore) SetLatestHeader(ctx context.Context, blockNo *big.Int, header *types.Header) error {
	log.Info().Str("module", "redis").Msgf("set latest %s", blockNo.String())

	err := s.client.Set(ctx, keyLatestBlockNo, blockNo.String(), 0).Err()
	if err != nil {
		log.Error().Str("module", "redis").Msgf(err.Error())
		return err
	}

	err = s.Store(ctx, blockNo, header)
	if err != nil {
		log.Error().Str("module", "redis").Msgf(err.Error())
		return err
	}

	return nil
}

func (s *redisStore) Latest(ctx context.Context) (*types.Header, error) {
	log.Info().Str("module", "redis").Msgf("latest")

	blockNoStr, err := s.client.Get(ctx, keyLatestBlockNo).Result()
	if err != nil {
		log.Error().Str("module", "redis").Msgf(err.Error())
		return nil, err
	}

	blockNo, ok := big.NewInt(0).SetString(blockNoStr, 10)
	if !ok {
		log.Error().Str("module", "redis").Msgf("cannot convert to big integer")
		return nil, errors.New("cannot convert to big integer")
	}

	header, err := s.Get(ctx, blockNo)
	if err != nil {
		log.Error().Str("module", "redis").Msgf(err.Error())
		return nil, err
	}

	return header, nil
}

func (s *redisStore) GetTotalDifficulty(ctx context.Context, blockNo *big.Int) (*big.Int, error) {
	log.Info().Str("module", "redis").Msgf("get total difficulty %s", blockNo.String())

	difficulty, err := s.client.Get(ctx, keyTotalDifficultyPrefix+blockNo.String()).Result()
	if err != nil {
		return nil, err
	}

	value, ok := big.NewInt(0).SetString(difficulty, 10)
	if !ok {
		return nil, errors.New("cannot convert to big integer")
	}

	return value, nil
}

func (s *redisStore) SetTotalDifficulty(ctx context.Context, blockNo *big.Int, totalDifficulty *big.Int) error {
	log.Info().Str("module", "redis").Msgf("set total difficulty %s %s", blockNo.String(), totalDifficulty.String())

	const expiration = time.Hour * 24 * 7

	fmt.Println(keyTotalDifficultyPrefix+blockNo.String(), totalDifficulty.String())
	err := s.client.Set(ctx, keyTotalDifficultyPrefix+blockNo.String(), totalDifficulty.String(), expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func newRedisClient(addr string) BlockHeaderStore {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &redisStore{client: rdb}
}
