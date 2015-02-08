package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

func Builds() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Build))
}

func BuildHistory(userId int64) []Build {
	var builds []Build
	Builds().Filter("UserId", userId).OrderBy("-CreateAt").All(&builds)
	return builds
}

func DeleteHistory(hid, uid int64) error {
	exist := Builds().Filter("Id", hid).Filter("UserId", uid).Exist()
	if exist {
		_, err := Builds().Filter("Id", hid).Delete()
		return err
	} else {
		return fmt.Errorf("no privilege")
	}
}
