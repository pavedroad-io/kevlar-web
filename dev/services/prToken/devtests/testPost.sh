#!/bin/bash

curl -H "Content-Type: application/json" \
     -X POST \
     -d @token.json \
     -v http://localhost:8081/api/v1/namespace/pavedroad.io/prTokens

