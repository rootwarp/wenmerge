#!/bin/bash

TARGET=53290823313153069154304

curl -X GET \
    -H "Content-Type: application/x-www-form-urlencoded" \
    localhost:9090/difficulty?target=$TARGET \
    -v
