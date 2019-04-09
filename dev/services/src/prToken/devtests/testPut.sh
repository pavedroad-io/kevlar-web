#!/bin/bash

curl -H "Content-Type: application/json" \
     -X PUT \
     -d @updatetoken.json \
     -v http://localhost:8081/api/v1/namespace/pavedroad.io/prTokens/ef3f33e8-42f8-4671-b15b-11ec4d61ad34

