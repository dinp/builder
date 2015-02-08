package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/dinp/builder/g"
	_ "github.com/go-sql-driver/mysql"
	systool "github.com/toolkits/sys"
	"html/template"
	"log"
	"os"
	"os/exec"
	"time"
)

type Build struct {
	Id       int64     `json:"id"`
	App      string    `json:"app" orm:"size(64)"`
	Version  string    `json:"version" orm:"size(64)"`
	Resume   string    `json:"resume" orm:"size(255)"`
	Base     string    `json:"base" orm:"size(1024)"`
	Image    string    `json:"image" orm:"size(1024)"`
	Tarball  string    `json:"tarball" orm:"size(1024)"`
	Repo     string    `json:"repo" orm:"size(255)"`
	Branch   string    `json:"branch" orm:"size(64)"`
	Status   string    `json:"status" orm:"size(255)"`
	UserId   int64     `json:"user_id" orm:"index"`
	UserName string    `json:"user_name" orm:"size(64)"`
	CreateAt time.Time `json:"create_at" orm:"auto_now_add;type(datetime)"`
}

func (this *Build) TableEngine() string {
	return "INNODB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci"
}

func init() {
	orm.RegisterModel(new(Build))
}

func (this *Build) UpdateBuildStatus(status string) (int64, error) {
	this.Status = status
	return orm.NewOrm().Update(this, "Status")
}

func (this *Build) UpdateImage(image string) (int64, error) {
	this.Image = image
	return orm.NewOrm().Update(this, "Image")
}

func (this *Build) GenDockerfile(workDir, tarball string) error {
	tpl, ok := g.TplMapping[this.Base]
	if !ok {
		return fmt.Errorf("no such base image: %s", this.Base)
	}

	t, err := template.ParseFiles(tpl)
	if err != nil {
		return err
	}

	out, err := os.Create(fmt.Sprintf("%s/Dockerfile", workDir))
	if err != nil {
		return err
	}

	defer out.Close()

	return t.Execute(out, map[string]interface{}{
		"Registry": g.Registry,
		"Tarball":  tarball,
		"AppDir":   "/opt/app",
	})
}

func (this *Build) DockerBuild(workDir string) {
	this.UpdateBuildStatus("docker building...")
	logFile := fmt.Sprintf("%s/%d.log", g.LogDir, this.Id)
	cmd := exec.Command(g.BuildScript, workDir, this.UserName, this.App, this.Version, g.Registry, logFile)
	err := cmd.Start()
	if err != nil {
		this.UpdateBuildStatus(fmt.Sprintf("error occur when build: %v", err))
		log.Printf("start cmd fail: %v when docker build %s/Dockerfile", err, workDir)
		return
	}

	err, timeout := systool.CmdRunWithTimeout(cmd, g.BuildTimeout)

	if err != nil {
		log.Printf("docker build %s/Dockerfile fail: %v", workDir, err)
	}

	if timeout {
		log.Printf("docker build %s/Dockerfile timeout", workDir)
	}

	if !timeout && err == nil {
		this.UpdateBuildStatus("successfully:-)")
		image := fmt.Sprintf("%s/%s/%s:%s", g.Registry, this.UserName, this.App, this.Version)
		this.UpdateImage(image)
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			defer f.Close()
			f.WriteString("\ndone:-)\n")
		}
		os.RemoveAll(workDir)
	} else {
		this.UpdateBuildStatus(fmt.Sprintf("is_timeout: %v, err: %v", timeout, err))
	}
}
