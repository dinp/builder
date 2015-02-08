package routers

import (
	"github.com/astaxie/beego"
	"github.com/dinp/builder/models"
	"github.com/toolkits/web"
	"strconv"
)

type Checker interface {
	CheckLogin()
}

type BaseController struct {
	beego.Controller
	CurrentUser *models.User
}

func (this *BaseController) Prepare() {
	if app, ok := this.AppController.(Checker); ok {
		app.CheckLogin()
	}
}

func (this *BaseController) SetPaginator(per int, nums int64) *web.Paginator {
	p := web.NewPaginator(this.Ctx.Request, per, nums)
	this.Data["paginator"] = p
	return p
}

func (this *BaseController) GetIntWithDefault(paramKey string, defaultVal int) int {
	valStr := this.GetString(paramKey)
	var val int
	if valStr == "" {
		val = defaultVal
	} else {
		var err error
		val, err = strconv.Atoi(valStr)
		if err != nil {
			val = defaultVal
		}
	}
	return val
}

func (this *BaseController) ServeErrJson(msg string) {
	this.Data["json"] = &models.ReturnDto{Msg: msg}
	this.ServeJson()
}

func (this *BaseController) ServeOKJson() {
	this.Data["json"] = &models.ReturnDto{Msg: ""}
	this.ServeJson()
}

func (this *BaseController) ServeDataJson(val interface{}) {
	this.Data["json"] = &models.ReturnDto{Msg: "", Data: val}
	this.ServeJson()
}
