package auth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"scheduler/db"
	"scheduler/log"
	"time"

	"github.com/satori/go.uuid"
)

const (
	_TOKEN  = "token"
	_NAME   = "name"
	_ACCESS = "access"
)

type AuthManager struct {
	log *logger.Log
	db  *db.Pgdb
}

func NewAuthManager(logDir string, dbParm *db.DbConnPara) *AuthManager {
	auth := new(AuthManager)

	auth.log = new(logger.Log)
	err := auth.log.InitLogger(logDir, "Auth")
	if nil != err {
		return nil
	}

	auth.db = db.NewPgDb(dbParm)
	if nil == auth.db {
		return nil
	}

	return auth
}

func (a *AuthManager) RegisterUser(username string, accessList []string) (string, error) {

	token := fmt.Sprintf("%s", uuid.NewV4())
	err := a.db.CreateNewUser(username, token, accessList)

	return token, err
}

func (a *AuthManager) getUserBaseToken(username string) string {
	userAuthInfo, err := a.db.GetUserAuthInfo(username)
	if nil != err {
		return ""
	}

	return userAuthInfo.Basetoken
}

func (a *AuthManager) getAccessList(username string) []string {
	userAuthInfo, err := a.db.GetUserAuthInfo(username)
	if nil != err {
		return []string{}
	}

	return userAuthInfo.AccList
}

func (a *AuthManager) DeleteUser(username string) error {
	err := a.db.DeleteUser(username)
	return err
}

func (a *AuthManager) CheckUserSignitural(username, timeStr, chksum string) (bool, error) {
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
