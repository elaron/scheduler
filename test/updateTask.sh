#!/bin/sh

curl 'http://127.0.0.1:2345/task?type=100&num=1'
curl -X PUT -H 'Content-Type: application/json' -d "{\"WorkerId\": \"worker001\", \"RequestId\": \"b943592a-b758-4844-b0a6-161022801361\", \"State\":2}" http://127.0.0.1:1234/task?type=100

