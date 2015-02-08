package g

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	filetool "github.com/toolkits/file"
	"log"
	"runtime"
	"strings"
	"time"
)

var (
	Debug        bool
	TmpDir       string
	LogDir       string
	BuildTimeout time.Duration
	UicInternal  string
	UicExternal  string
	Registry     string
	Cache        cache.Cache
	BuildScript  string
	TplMapping   map[string]string = make(map[string]string)
	Token        string
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	ParseConfig()
}

func ParseConfig() {

	runMode := beego.AppConfig.String("runmode")
	if runMode == "dev" {
		Debug = true
	} else {
		Debug = false
	}

	TmpDir = beego.AppConfig.String("tmpdir")
	if TmpDir == "" {
		log.Fatalln("configuration tmpdir is blank")
	}

	err := filetool.InsureDir(TmpDir)
	if err != nil {
		log.Fatalf("create dir: %s fail: %v", TmpDir, err)
	}

	TmpDir, err = filetool.RealPath(TmpDir)
	if err != nil {
		log.Fatalf("get real path of %s fail: %v", TmpDir, err)
	}

	LogDir = beego.AppConfig.String("logdir")
	if LogDir == "" {
		log.Fatalln("configuration logdir is blank")
	}

	err = filetool.InsureDir(LogDir)
	if err != nil {
		log.Fatalf("create dir: %s fail: %v", LogDir, err)
	}

	LogDir, err = filetool.RealPath(LogDir)
	if err != nil {
		log.Fatalf("get real path of %s fail: %v", LogDir, err)
	}

	Token = beego.AppConfig.String("token")

	UicInternal = beego.AppConfig.String("uicinternal")
	if UicInternal == "" {
		log.Fatalln("configuration uicinternal is blank")
	}

	UicExternal = beego.AppConfig.String("uicexternal")
	if UicExternal == "" {
		log.Fatalln("configuration uicexternal is blank")
	}

	BuildScript = beego.AppConfig.String("buildscript")
	if BuildScript == "" {
		log.Fatalln("configuration buildscript is blank")
	}

	Registry = beego.AppConfig.String("registry")
	if Registry == "" {
		log.Fatalln("configuration registry is blank")
	}

	_buildTimeout, err := beego.AppConfig.Int64("buildtimeout")
	if err != nil {
		log.Fatalf("parse configuration buildtimeout fail: %v", err)
	}

	BuildTimeout = time.Duration(_buildTimeout) * time.Minute

	// tpl mapping
	tpl_mapping := beego.AppConfig.String("tplmapping")
	tpl_mapping = strings.TrimSpace(tpl_mapping)
	mappings := strings.Split(tpl_mapping, ",")
	for i := 0; i < len(mappings); i++ {
		_mappings := strings.TrimSpace(mappings[i])
		kv := strings.Split(_mappings, "=>")
		if len(kv) != 2 {
			log.Fatalf("split %s fail", _mappings)
		}

		TplMapping[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	// cache
	Cache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		log.Fatalln("start cache fail :-(")
	}

	// db
	dbuser := beego.AppConfig.String("dbuser")
	dbpass := beego.AppConfig.String("dbpass")
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbname := beego.AppConfig.String("dbname")
	dblink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbuser, dbpass, dbhost, dbport, dbname)
	// dblink = "root:1234@/uic?charset=utf8&loc=Asia%2FChongqing"

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", dblink+"&loc=Asia%2FChongqing", 30, 200)
	// orm.DefaultTimeLoc = time.UTC

	if Debug {
		orm.Debug = true
	}
}
