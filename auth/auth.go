package auth

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/satori/go.uuid"
	"math"
	"scheduler/common"
	"scheduler/log"
	"time"
)

const (
	_TOKEN  = "token"
	_NAME   = "name"
	_ACCESS = "access"
)

type Auth struct {
	pool *pool.Pool
	log  logger.Log
}

func (a *Auth) Init(address string, port int, poolSize int, prefix string) {
	addr := fmt.Sprintf("%s:%d", address, port)
	p, err := pool.New("tcp", addr, poolSize)
	if nil != err {
		a.log.Info.Println("Connect to redis fail, ", err)
		return
	}

	a.pool = p

	tag := fmt.Sprintf("%s-auth-", prefix)
	a.log.InitLogger(tag)
}

func (a *Auth) RegisterUser(username string, accessList []int) (string, error) {

	var token string
	field := comm.GetUserTable(username)

	resp := a.pool.Cmd("HEXISTS", field, _TOKEN)
	if nil == resp.Err {
		exist, _ := resp.Int()
		if 1 == exist {
			resp = a.pool.Cmd("HGET", field, _TOKEN)
			if nil == resp.Err {
				token, _ = resp.Str()
			}
		}
	}

	if len(token) == 0 {
		token = fmt.Sprintf("%s", uuid.NewV4())
	}

	access, err := json.Marshal(accessList)
	if nil != err {
		s := fmt.Sprintf("Encoding accessList fail, %s", err.Error())
		a.log.Info.Println(s)
		return "", errors.New(s)
	}

	c, err := a.pool.Get()
	if nil != err {
		s := fmt.Sprintf("Get connection from pool fail, ", err.Error())
		a.log.Info.Println(s)
		return "", errors.New(s)
	}
	defer a.pool.Put(c)

	c.Cmd("MULTI")
	c.Cmd("HSET", field, _NAME, username)
	c.Cmd("HSET", field, _TOKEN, token)
	c.Cmd("HSET", field, _ACCESS, access)
	c.Cmd("EXEC")

	return token, nil
}

func (a *Auth) getUserBaseToken(username string) string {
	field := comm.GetUserTable(username)
	resp := a.pool.Cmd("HGET", field, _TOKEN)
	if nil != resp.Err {
		a.log.Info.Println("Get user base token fail, ", username, resp.Err)
		return ""
	}

	token, err := resp.Str()
	if nil != err {
		a.log.Info.Println("Decode base token fail, ", username, err)
		return ""
	}

	return token
}

func (a *Auth) getAccessList(username string) []int {
	field := comm.GetUserTable(username)
	resp := a.pool.Cmd("HGET", field, _ACCESS)
	if nil != resp.Err {
		a.log.Info.Println("Get user base access fail, ", username, resp.Err)
		return []int{}
	}

	access, err := resp.Str()
	if nil != err {
		a.log.Info.Println("Decode base access fail, ", username, err)
		return []int{}
	}

	var accessList []int
	err = json.Unmarshal([]byte(access), &accessList)
	if nil != err {
		a.log.Info.Println("Decode access list from json fail, ", username, err)
		return []int{}
	}

	return accessList
}

func (a *Auth) DeleteUser(username string) {
	field := comm.GetUserTable(username)
	a.pool.Cmd("DEL", field)
}

func (a *Auth) CheckUserSignitural(username, timeStr, chksum string) (bool, error) {
	token := a.getUserBaseToken(username)
	if len(token) == 0 {
		s := fmt.Sprintf("Get %s base token fail.", username)
		a.log.Info.Println(s)
		return false, errors.New(s)
	}

	if isTimeOk(timeStr, time.Now()) == false {
		s := fmt.Sprintf("Illegal timestamp")
		a.log.Info.Println(s, timeStr, time.Now())
		return false, errors.New(s)
	}
	return identify(username, token, timeStr, chksum), nil
}

func isTimeOk(timeStr string, curr time.Time) bool {

	t1, err := time.Parse("2006-01-02 03:04:05", timeStr)
	if nil != err {
		return false
	}

	deltaTime := math.Abs(float64(curr.Sub(t1).Seconds()))
	fmt.Println("delta time:", t1, curr, deltaTime)
	if deltaTime > 300 {
		return false
	}

	return true
}

func identify(username, baseToken, timeStr, chksum string) bool {

	str := fmt.Sprintf("%s%s%s", username, baseToken, timeStr)
	checksum := sha256.Sum256([]byte(str))

	sum := fmt.Sprintf("%x", checksum)
	if sum != chksum {
		return false
	} else {
		return true
	}
}
