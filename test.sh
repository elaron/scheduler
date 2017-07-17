#!/bin/sh

curl http://127.0.0.1:1234/request?type=100
curl -X POST -d "{\"test\": \"that\"}" http://127.0.0.1:1234/request?type=200
curl -X PUT -d "{\"test\": \"that\"}" http://127.0.0.1:1234/request?type=300
curl -X DELETE -d "{\"test\": \"that\"}" http://127.0.0.1:1234/request?type=400
