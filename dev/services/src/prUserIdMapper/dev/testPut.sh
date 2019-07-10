#!/bin/bash

curl -H "Content-Type: application/json" \
     -X PUT \
     -d @userIdMapPutData.json \
     -v http://localhost:8082/api/v1/namespace/pavedroad.io/prUserIdMappers/john@scharber.com

