package routers

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/dinp/builder/g"
	"github.com/dinp/builder/models"
	filetool "github.com/toolkits/file"
	str_ "github.com/toolkits/str"
	"regexp"
	"strconv"
	"strings"
)

type MainController struct {
	PageController
}

func (this *MainController) Get() {
	this.Data["Mappings"] = g.TplMapping
	this.Layout = "layout.html"
	this.TplNames = "index.html"
}

func (this *MainController) Build() {
	o := &models.Build{}
	o.App = this.GetString("app")
	o.Version = this.GetString("version")
	o.Resume = this.GetString("resume")
	o.Base = this.GetString("base")
	o.Tarball = this.GetString("tarball")

	defer func() {
		this.Data["Mappings"] = g.TplMapping
		this.Data["O"] = o
		this.Layout = "layout.html"
		this.TplNames = "index.html"
	}()

	if o.App == "" {
		this.Data["Msg"] = "app名称不能为空"
		return
	}

	// app should be a-zA-Z0-9_-
	var appPattern = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9\-\_]*$`)
	if !appPattern.MatchString(o.App) {
		this.Data["Msg"] = "app名称应该符合正则：/^[a-zA-Z_]+[a-zA-Z0-9\\-\\_]*$/"
		return
	}

	if o.Version == "" {
		this.Data["Msg"] = "版本不能为空"
		return
	}

	// version should be digit and .
	var versionPattern = regexp.MustCompile(`^[0-9]+[0-9\.]*$`)
	if !versionPattern.MatchString(o.Version) {
		this.Data["Msg"] = "version应该符合正则：/^[0-9]+[0-9\\.]*$/"
		return
	}

	if o.Base == "" {
		this.Data["Msg"] = "Base Image不能为空"
		return
	}

	workDir := fmt.Sprintf("%s/%s", g.TmpDir, str_.RandSeq(6))
	err := filetool.InsureDir(workDir)
	if err != nil {
		this.Data["Msg"] = fmt.Sprintf("create temp dir fail: %v", err)
		return
	}

	var fileName, filePath string

	if o.Tarball == "" {
		_, header, err := this.GetFile("file")
		if err != nil {
			this.Data["Msg"] = err.Error()
			return
		}

		// handle upload file
		fileName = header.Filename
		filePath = fmt.Sprintf("%s/%s", workDir, fileName)
		err = this.SaveToFile("file", filePath)
		if err != nil {
			this.Data["Msg"] = fmt.Sprintf("save file fail: %v", err)
			return
		}
	} else {
		if !strings.HasPrefix(o.Tarball, "http://") {
			this.Data["Msg"] = "tarball地址应该是一个http地址"
			return
		}

		if !(strings.HasSuffix(o.Tarball, ".tar.gz") || strings.HasSuffix(o.Tarball, ".war")) {
			this.Data["Msg"] = "tarball地址应该以.tar.gz或.war结尾"
			return
		}

		idx := strings.LastIndex(o.Tarball, "/")
		fileName = o.Tarball[idx+1:]
		filePath = fmt.Sprintf("%s/%s", workDir, fileName)

		err = filetool.Download(filePath, o.Tarball)
		if err != nil {
			this.Data["Msg"] = fmt.Sprintf("download tarball fail: %v", err)
			return
		}
	}

	o.UserId = this.CurrentUser.Id
	o.UserName = this.CurrentUser.Name
	o.Status = "saved meta in db"
	_, err = orm.NewOrm().Insert(o)
	if err != nil {
		this.Data["Msg"] = fmt.Sprintf("save meta to db fail: %v", err)
		return
	}

	err = o.GenDockerfile(workDir, fileName)
	if err != nil {
		this.Data["Msg"] = fmt.Sprintf("generate Dockerfile fail: %v", err)
		return
	}

	go o.DockerBuild(workDir)

	this.Redirect(fmt.Sprintf("/progress/%d", o.Id), 302)
}

func (this *MainController) Progress() {
	idStr := this.Ctx.Input.Param(":id")
	this.Data["Id"] = idStr

	buildId, err := strconv.ParseInt(idStr, 10, 64)
	build := models.Build{Id: buildId}
	err = orm.NewOrm().Read(&build)
	if err != nil {
		this.ServeErrJson(err.Error())
		return
	}

	this.Data["Image"] = fmt.Sprintf("%s/%s/%s:%s", g.Registry, build.UserName, build.App, build.Version)
	this.Layout = "layout.html"
	this.TplNames = "progress.html"
}

func (this *MainController) Log() {
	idStr := this.Ctx.Input.Param(":id")
	buildId, err := strconv.ParseInt(idStr, 10, 64)
	build := models.Build{Id: buildId}
	err = orm.NewOrm().Read(&build)
	if err != nil {
		this.ServeErrJson(err.Error())
		return
	}

	content, err := filetool.ToString(fmt.Sprintf("%s/%d.log", g.LogDir, buildId))
	if err != nil {
		this.ServeErrJson(err.Error())
		return
	}

	content = strings.Replace(content, "\n", "<br>", -1)

	this.ServeDataJson(map[string]interface{}{"build": build, "log": content})
}

func (this *MainController) History() {
	this.Data["BuildHistory"] = models.BuildHistory(this.CurrentUser.Id)
	this.Layout = "layout.html"
	this.TplNames = "history.html"
}

func (this *MainController) DeleteHistory() {
	id, err := this.GetInt("id")
	if err != nil {
		this.ServeErrJson("id invalid")
		return
	}

	err = models.DeleteHistory(int64(id), this.CurrentUser.Id)
	if err != nil {
		this.ServeErrJson(err.Error())
	} else {
		this.ServeOKJson()
	}
}
