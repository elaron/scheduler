#!/bin/sh

curl http://127.0.0.1:1234/request?type=100
curl -X  POST -H 'Content-Type: application/json' -d "{\"test\": \"a new request\"}" http://127.0.0.1:1234/request?type=100
curl -X PUT -H 'Content-Type: application/json' -d "{\"WorkerId\": \"worker001\", \"RequestId\": \"11121212\", \"State\":2}" http://127.0.0.1:1234/request?type=100
curl -X DELETE -d "{\"reqid\": \"11121212\"}" http://127.0.0.1:1234/request?type=100

curl http://127.0.0.1:1234/clean?type=100
