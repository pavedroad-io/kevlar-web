#!/bin/bash

curl -X DELETE \
     -H "Content-Type: application/json" \
     -v http://localhost:8082/api/v1/namespace/pavedroad.io/prUserIdMappers/john@scharber.com


