#!/bin/bash

TARGET=52863991973008772372162

curl -X GET \
    -H "Content-Type: application/x-www-form-urlencoded" \
    localhost:9090/difficulty?target=$TARGET \
    -v
