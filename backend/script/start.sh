#!/bin/bash


docker run --rm \
    -p 9090:9090 \
    --env ETH_WSS_URL=wss://eth-mainnet.alchemyapi.io/v2/cFAw9HT7YQUSRKHktINfqL5_hkQALWYV \
    --env ETH_RPC_URL=https://eth-mainnet.alchemyapi.io/v2/cFAw9HT7YQUSRKHktINfqL5_hkQALWYV \
    rootwarp/wenmerge-backend:latest
