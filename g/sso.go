package g

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
)

func GetSig() (sig string, err error) {
	uri := fmt.Sprintf("%s/sso/sig", UicInternal)
	req := httplib.Get(uri)
	sig, err = req.String()
	return
}

func Logout(sig string) error {
	uri := fmt.Sprintf("%s/sso/logout/%s", UicInternal, sig)
	req := httplib.Get(uri)
	req.Param("token", Token)
	_, err := req.String()
	return err
}
