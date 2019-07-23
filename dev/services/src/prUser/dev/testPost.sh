#!/bin/bash

curl -H "Content-Type: application/json" \
     -X POST \
     -d @user.json \
     -v http://localhost:8083/api/v1/namespace/pavedroad.io/prUsers

