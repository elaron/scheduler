package auth

import (
	"testing"
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
