package main

import (
	"github.com/astaxie/beego"
	_ "github.com/dinp/builder/g"
	_ "github.com/dinp/builder/routers"
)

func main() {
	beego.Run()
}
