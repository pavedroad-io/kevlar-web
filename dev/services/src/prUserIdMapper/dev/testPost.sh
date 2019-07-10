#!/bin/bash

curl -H "Content-Type: application/json" \
     -X POST \
     -d @userIdMap.json \
     -v http://localhost:8082/api/v1/namespace/pavedroad.io/prUserIdMappers

