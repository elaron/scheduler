package auth

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/satori/go.uuid"
)

const (
	_TOKEN  = "token"
	_NAME   = "name"
	_ACCESS = "access"
)

func registerUser(username string, accessList []int) (string, error) {

	token := fmt.Sprintf("%s", uuid.NewV4())
	field := getUserTable(username)

	access, err := json.Marshal(accessList)
	if nil != err {
		s := fmt.Sprintf("Encoding accessList fail, %s", err.Error())
		g_log.Info.Println(s)
		return "", errors.New(s)
	}

	c, err := g_redisPool.Get()
	if nil != err {
		s := fmt.Sprintf("Get connection from pool fail, ", err.Error())
		g_log.Info.Println(s)
		return "", errors.New(s)
	}
	defer g_redisPool.Put(c)

	c.Cmd("MULTI")
	c.Cmd("HSET", field, _NAME, username)
	c.Cmd("HSET", field, _TOKEN, token)
	c.Cmd("HSET", field, _ACCESS, access)
	c.Cmd("EXEC")

	return token, nil
}

func getUserBaseToken(username string) string {
	field := getUserTable(username)
	resp := g_redisPool.Cmd("HGET", field, _TOKEN)
	if nil != resp.Err {
		g_log.Info.Println("Get user base token fail, ", username, resp.Err)
		return ""
	}

	token, err := resp.Str()
	if nil != err {
		g_log.Info.Println("Decode base token fail, ", username, err)
		return ""
	}

	return token
}

func getAccessList(username string) []int {
	field := getUserTable(username)
	resp := g_redisPool.Cmd("HGET", field, _ACCESS)
	if nil != resp.Err {
		g_log.Info.Println("Get user base access fail, ", username, resp.Err)
		return ""
	}

	access, err := resp.Str()
	if nil != err {
		g_log.Info.Println("Decode base access fail, ", username, err)
		return ""
	}

	var accessList []int
	err = json.Unmarshal([]byte(access), &accessList)
	if nil != err {
		g_log.Info.Println("Decode access list from json fail, ", username, err)
		return ""
	}

	return accessList
}

func deleteUser(username string) {
	field := getUserTable(username)
	g_redisPool.Cmd("DEL", field)
}
