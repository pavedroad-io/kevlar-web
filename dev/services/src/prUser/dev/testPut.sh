#!/bin/bash

export uuid=`curl -H "Content-Type: application/json" -s http://localhost:8083/api/v1/namespace/pavedroad.io/prUsers/fool?query=name | jq -r '.Metadata.uuid'`

echo "UUID for user test is :  $uuid"

curl -H "Content-Type: application/json" \
     -X PUT \
     -d @userPutData.json \
     -v http://localhost:8083/api/v1/namespace/pavedroad.io/prUsers/$uuid

