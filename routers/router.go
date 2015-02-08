package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &MainController{})
	beego.Router("/health", &ApiController{}, "get:Health")
	beego.Router("/logout", &ApiController{}, "get:Logout")
	beego.Router("/build", &MainController{}, "post:Build")
	beego.Router("/progress/:id:int", &MainController{}, "get:Progress")
	beego.Router("/history", &MainController{}, "get:History;post:DeleteHistory")
	beego.Router("/log/:id:int", &MainController{}, "get:Log")
}
