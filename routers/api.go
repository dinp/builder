package routers

import (
	"github.com/dinp/builder/g"
)

type ApiController struct {
	BaseController
}

func (this *ApiController) Health() {
	this.Ctx.WriteString("ok")
}

func (this *ApiController) Logout() {
	sig := this.Ctx.GetCookie("sig")
	if sig == "" {
		this.Ctx.WriteString("u'r not login")
		return
	}

	err := g.Logout(sig)
	if err != nil {
		this.Ctx.WriteString("logout from uic fail")
		return
	}

	this.Ctx.SetCookie("sig", "", 0, "/")
	this.Redirect("/", 302)
}
