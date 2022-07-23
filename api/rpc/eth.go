package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
)

var (
	rpcURL = os.Getenv("ETH_RPC_URL")
)

type httpBlock struct {
	JSONRpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Hash            string `json:"hash"`
		Number          string `json:"number"`
		Difficulty      string `json:"difficulty"`
		TotalDifficulty string `json:"totalDifficulty"`
	} `json:"result"`
}

// Block is ...
type Block struct {
	Hash            string
	Number          *big.Int
	Difficulty      *big.Int
	TotalDifficulty *big.Int
}

// Client provices interfaces for eth client.
type Client interface {
	GetBlockNumber() (*big.Int, error)
	GetBlockByNumber(blockNo *big.Int) (*Block, error)
}

type ethClient struct {
}

func (c *ethClient) GetBlockNumber() (*big.Int, error) {
	cli := http.Client{}

	reqBody := struct {
		Jsonrpc string   `json:"jsonrpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
		ID      int      `json:"id"`
	}{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		ID:      1,
	}

	rawBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, rpcURL, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)

	respData := struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  string `json:"result"`
	}{}

	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return nil, err
	}

	i := new(big.Int)
	i, ok := i.SetString(strings.TrimPrefix(respData.Result, "0x"), 16)
	if !ok {
		return nil, errors.New("invalid hex number")
	}

	return i, nil
}

func (c *ethClient) GetBlockByNumber(blockNo *big.Int) (*Block, error) {
	cli := http.Client{}

	reqBody := struct {
		Jsonrpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		ID      int           `json:"id"`
	}{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		ID:      1,
		Params: []interface{}{
			"0x" + blockNo.Text(16),
			true,
		},
	}

	rawBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, rpcURL, bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)

	b := httpBlock{}
	err = json.Unmarshal(respBody, &b)
	if err != nil {
		return nil, err
	}

	retBlock := Block{Hash: b.Result.Hash}

	blockNo, ok := new(big.Int).SetString(strings.TrimPrefix(b.Result.Number, "0x"), 16)
	if !ok {
		return nil, errors.New("invalid hex format")
	}

	difficulty, ok := new(big.Int).SetString(strings.TrimPrefix(b.Result.Difficulty, "0x"), 16)
	if !ok {
		return nil, errors.New("invalid hex format")
	}

	totalDifficulty, ok := new(big.Int).
		SetString(strings.TrimPrefix(b.Result.TotalDifficulty, "0x"), 16)
	if !ok {
		return nil, errors.New("invalid hex format")
	}

	retBlock.Number = blockNo
	retBlock.Difficulty = difficulty
	retBlock.TotalDifficulty = totalDifficulty

	return &retBlock, nil
}

// NewClient create new client instance.
func NewClient() Client {
	return &ethClient{}
}
