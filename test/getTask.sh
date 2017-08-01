#!/bin/sh

num=$1
curl "http://127.0.0.1:6668/task?type=100&num=${num}"

