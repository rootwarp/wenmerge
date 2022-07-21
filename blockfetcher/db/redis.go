package db

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v9"
)

type redisStore struct {
	client *redis.Client
}

func (s *redisStore) Store(ctx context.Context, blockNo *big.Int, header *types.Header) error {
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

func (s *redisStore) Latest(ctx context.Context) (*types.Header, error) {
	// TODO:
	return nil, nil
}

func newRedisClient(addr string) BlockHeaderStore {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &redisStore{client: rdb}
}
