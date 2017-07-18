#!/bin/sh

num=$1
curl "http://127.0.0.1:2345/task?type=100&num=${num}"

