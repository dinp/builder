package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/dinp/builder/g"
)

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Im    string `json:"im"`
	Phone string `json:"phone"`
}

func GetUser(sig string) (*User, error) {
	key := fmt.Sprintf("u:%s", sig)
	u := g.Cache.Get(key)
	if u != nil {
		uobj := u.(User)
		return &uobj, nil
	}

	uri := fmt.Sprintf("%s/sso/user/%s", g.UicInternal, sig)
	req := httplib.Get(uri)
	req.Param("token", g.Token)
	resp, err := req.Response()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	type TmpStruct struct {
		User *User `json:"user"`
	}
	var t TmpStruct
	err = decoder.Decode(&t)
	if err != nil {
		return nil, err
	}

	// don't worry cache expired. we just use username which can not modify
	g.Cache.Put(key, *t.User, int64(360000))

	return t.User, nil
}
