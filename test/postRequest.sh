#!/bin/sh

str=$1
curl -X  POST -H 'Content-Type: application/json' -d "{\"test\": \"a new request--${str}\"}" http://127.0.0.1:1234/request?type=100
