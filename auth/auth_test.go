package auth

import (
	"testing"
	"time"
)

func Test_RegisterUser_1(t *testing.T) {
	var a Auth
	a.Init("localhost", 6379, 10, "scheduler")

	username := "aaa"
	accessList := []int{100, 101}
	a.DeleteUser(username)

	token1, err := a.RegisterUser(username, accessList)
	if nil != err {
		t.Error("CASE: [Register user] fail!")
		return
	}

	token2, err := a.RegisterUser(username, accessList)
	if nil != err {
		t.Error("CASE: [Register user] Fail! Cannot register repeately.")
		return
	}

	if token2 != token1 {
		t.Error("CASE: [Register user] Fail! Get different token when registing repeately. ")
		return
	}

	t.Log("CASE [Register user] pass!")
}

func Test_Identify_1(t *testing.T) {
	username := "bbb"
	baseToken := "b84f513c-fd3e-418b-b4f7-575076754314"
	timeStr := "2017-07-25 10:00:00"
	checksum := "6eda021e1f8dc5e9f5c1a69f2fd0290ae6337d7f928c8366fd13aa2600dfdbd6"
	ok := identify(username, baseToken, timeStr, checksum)
	if true != ok {
		t.Error("[CASE: identify check fail]")
	} else {
		t.Log("[CASE: identify check success]")
	}
}

func Test_IsTimeOk_1(t *testing.T) {
	beTest := "2017-07-25 10:00:00"
	t1, _ := time.Parse("2006-01-02 03:04:05", "2017-07-25 09:55:00")
	t2, _ := time.Parse("2006-01-02 03:04:05", "2017-07-25 10:05:00")
	t3, _ := time.Parse("2006-01-02 03:04:05", "2017-07-25 09:54:59")
	t4, _ := time.Parse("2006-01-02 03:04:05", "2017-07-25 10:05:01")

	ok := isTimeOk(beTest, t1)
	if ok != true {
		t.Error("[CASE: isTimeOk check fail 1]")
	}

	ok = isTimeOk(beTest, t2)
	if ok != true {
		t.Error("[CASE: isTimeOk check fail 2]")
	}

	ok = isTimeOk(beTest, t3)
	if ok != false {
		t.Error("[CASE: isTimeOk check early fail]")
	}

	ok = isTimeOk(beTest, t4)
	if ok != false {
		t.Error("[CASE: isTimeOk check later fail]")
	}
}
