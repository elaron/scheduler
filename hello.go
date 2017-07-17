package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"time"
)

func addRequest(c redis.Conn, request string) {

	ts := time.Now().UnixNano()
	id := uuid.NewV4()
	fmt.Println(ts, id)

	v, err := c.Do("ZADD", "request", ts, id)
	if nil != err {
		fmt.Println("ZADD request fail", err)
		return
	}

	fmt.Println(v)
	values, err := redis.Values(c.Do("ZRANGE", "request", "0", "-1", "WITHSCORES"))
	if nil != err {
		fmt.Println("Get requst list fail, ", err)
		return
	} else {
		fmt.Println("Request list:")
		for k, req := range values {
			if k%2 == 0 {
				fmt.Printf("id: %d %s ", k, req)
			} else {
				fmt.Printf("score: %s\n", req)
			}
		}
	}
}

func cleanRequestList(c redis.Conn, name string) {
	_, err := c.Do("ZREMRANGEBYRANK", name, "0", "-1")
	if nil != err {
		fmt.Printf("Clean %s fail, %s\n", name, err.Error())
		return
	}
}

func main() {

	c, err := redis.Dial("tcp", "localhost:6379")
	if nil != err {
		fmt.Println("Connect to redis fail, ", err)
		return
	}
	defer c.Close()
	cleanRequestList(c, "request")

	for i := 0; i < 5; i++ {
		addRequest(c, "Elar")
	}
}
