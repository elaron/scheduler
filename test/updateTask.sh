#!/bin/sh

reqID=$1

curl -X PUT -H 'Content-Type: application/json' -d "{\"WorkerId\": \"worker001\", \"RequestId\": \"${reqID}\", \"State\":2}" http://127.0.0.1:6666/task?type=100

