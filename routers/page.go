package routers

import (
	"fmt"
	"github.com/dinp/builder/g"
	"github.com/dinp/builder/models"
)

type PageController struct {
	BaseController
}

func (this *PageController) hasBusinessCookie() bool {
	if this.CurrentUser != nil {
		return true
	}

	sigInCookie := this.Ctx.GetCookie("sig")
	if sigInCookie != "" {
		u, err := models.GetUser(sigInCookie)
		if err != nil || u == nil {
			return false
		}

		this.CurrentUser = u
		return true
	}

	return false
}

func (this *PageController) CheckLogin() {
	if this.hasBusinessCookie() {
		this.Data["CurrentUser"] = this.CurrentUser
		return
	}

	sig, err := g.GetSig()
	if err != nil {
		this.ServeErrJson(fmt.Sprintf("curl get sig fail: %v", err))
		return
	}

	this.Ctx.SetCookie("sig", sig, 2592000, "/")
	this.Redirect(fmt.Sprintf("%s/auth/login?sig=%s&callback=http://%s:%d%s", g.UicExternal, sig, this.Ctx.Input.Host(), this.Ctx.Input.Port(), this.Ctx.Input.Uri()), 302)
}
