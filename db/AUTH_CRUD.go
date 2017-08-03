package db

import (
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

//create table user_auth_info (username varchar(64) PRIMARY KEY, basetoken varchar(64) NOT NULL, accesslist varchar(1024));
const (
	_USER_AUTH_INFO_TABLENAME_ = "user_auth_info"
	_FIELD_USERNAME_           = "username"
	_FIELD_BASETOKEN_          = "basetoken"
	_FIELD_ACCESSLIST_         = "accesslist"
)

type UserAuthInfo struct {
	Username  string
	Basetoken string
	AccList   []string
}

func (p *Pgdb) CreateNewUser(username, baseToken string, accessList []string) error {

	db := p.db
	accList, _ := json.Marshal(accessList)

	cmd := fmt.Sprintf("insert into %s(%s,%s,%s) values('%s','%s','%s');",
		_USER_AUTH_INFO_TABLENAME_,
		_FIELD_USERNAME_, _FIELD_BASETOKEN_, _FIELD_ACCESSLIST_,
		username, baseToken, string(accList))

	_, err := db.Exec(cmd)
	if nil != err {
		return err
	}
	return nil
}

func (p *Pgdb) DeleteUser(username string) error {

	db := p.db

	cmd := fmt.Sprintf("delete from %s where %s = '%s';",
		_USER_AUTH_INFO_TABLENAME_,
		_FIELD_USERNAME_,
		username)

	_, err := db.Exec(cmd)
	if nil != err {
		return err
	}
	return nil
}

func (p *Pgdb) GetUserAuthInfo(username string) (UserAuthInfo, error) {

	var name, basetoken, acclist string
	db := p.db

	cmd := fmt.Sprintf("select * from %s where %s = '%s';",
		_USER_AUTH_INFO_TABLENAME_,
		_FIELD_USERNAME_,
		username)

	rows, err := db.Query(cmd)
	if nil != err {
		goto Fail
	}

	for rows.Next() {
		err = rows.Scan(&name, &basetoken, &acclist)
		if nil != err {
			goto Fail
		}

		aui := UserAuthInfo{Username: name, Basetoken: basetoken}
		err = json.Unmarshal([]byte(acclist), &aui.AccList)
		if nil != err {
			goto Fail
		}

		return aui, nil
	}

Fail:
	return UserAuthInfo{}, err
}
