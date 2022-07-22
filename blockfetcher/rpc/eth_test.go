package rpc

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlockNumber(t *testing.T) {
	cli := NewClient()
	block, err := cli.GetBlockNumber()

	assert.Nil(t, err)
	assert.NotNil(t, block)
}

/*
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "number": "0x186a0",
    "difficulty": "0x37f5c0b3e61",
    "extraData": "0x476574682f76312e302e312d38326566323666362f6c696e75782f676f312e34",
    "gasLimit": "0x2fefd8",
    "gasUsed": "0x0",
    "hash": "0x91c90676cab257a59cd956d7cb0bceb9b1a71d79755c23c7277a0697ccfaf8c4",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "miner": "0xe6a7a1d47ff21b6321162aea7c6cb457d5476bca",
    "mixHash": "0x9769104bc21358e422d8b0938334c16072e5168cfdbba49614e4cb821ff26176",
    "nonce": "0x9f63aafeec219854",
    "parentHash": "0xfbafb4b7b6f6789338d15ff046f40dc608a42b1a33b093e109c6d7a36cd76f61",
    "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "size": "0x225",
    "stateRoot": "0x209230089ff328b2d87b721c48dbede5fd163c3fae29920188a7118275ab2013",
    "timestamp": "0x55d19762",
    "totalDifficulty": "0x259ec16a4e5ec78",
    "transactions": [],
    "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "uncles": []
  }
}
*/
func TestGetBlockByNumber(t *testing.T) {
	cli := NewClient()
	b, err := cli.GetBlockByNumber(big.NewInt(100000))

	assert.Nil(t, err)

	assert.Equal(t, "0x91c90676cab257a59cd956d7cb0bceb9b1a71d79755c23c7277a0697ccfaf8c4", b.Hash)
	assert.Equal(t, big.NewInt(100000), b.Number)
	assert.Equal(t, b.Difficulty.Text(16), "37f5c0b3e61")
	assert.Equal(t, b.TotalDifficulty.Text(16), "259ec16a4e5ec78")
}
