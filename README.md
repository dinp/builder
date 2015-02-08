Builder
=======

这是DINP中的一个编配平台，核心作用就是把用户提供的代码包，变成一个image推送到Docker-Registry

# 设计理念

- 平台不管编译，用户自己编译好，然后把编译好的代码扔给平台
- 不管是什么语言的代码，最终都打包成一个docker image，通过docker来做规范化

## 编译与否的问题

用户的代码要部署在PaaS平台上，首先要做的就是把代码扔给平台，那么问题来了，通过什么方式扔给平台呢？通过git？还是直接上传？扔上来的是源代码？还是已经编译完成的？解释型语言当然是不需要编译的，但是编译型的呢，比如golang，java，我们要让用户给我们.go、.java文件然后我们去编译？还是用户编译好，直接给我们二进制，给我们.class文件？

平台刚起步，我们希望在可接受的范围内，做的事情越少越好，这样出错的几率就越少，之后平台慢慢发展，可以再加入一些新的feature。所以我们选择让用户把编译之后的代码交给我们，比如java的话，就扔个war包给我们即可；golang的话需要编译成二进制（编译成64位Linux的，因为平台是64位Linux）。

有些PaaS，比如tsuru，是让用户把代码push到某个指定repo的指定分支的，然后后端的git receiver就会触发编译、上线脚本。看起来是挺好的，但是这样处理比较麻烦，坑太多，以后再说吧

接下来是要确定通过什么方式把代码扔给平台，我们目前支持两种方式，一个是在页面直接上传，一个是给一个可以下载的http地址，平台去下载。

## 使用docker来做规范化打包

用户的代码可能是不同语言写的，比较常见的比如Java、PHP、Python，用户给我们的是一个tar.gz包，我们要怎么把其变成一个docker image呢？

最简单的实现方式：dinp提供一个base image，里边就只是Ubuntu或者centos的环境，用户自己搞定runtime依赖，比如Java的话，用户的tar.gz包中应该包含JDK、tomcat、webapp，然后约定一个启动脚本，比如就是根目录的control文件，平台只需要`./control start`即可启动。如果用户是PHP的代码，需要把nginx、php-fpm、code打包进来。

但是这种方式太麻烦，如果是golang的话可以使用上面的方式，因为golang是静态编译的，也没啥依赖。Java、php、Python、Ruby让用户这么搞，那不得哭了……

我们可以针对不同的语言做不同的base image，比如Java的，base image中提前把JDK和tomcat做好，用户直接上传一个war包，搞定！比如PHP的base image，提前把nginx、php-fpm之类的做好，用户直接把php code上传，搞定！

嗯，这也是Builder平台现在采用的方式，我们做了一些base image，让用户在页面上选择使用哪个，用户选择了base image，也就间接的告诉了我们他是什么类型的代码了（比如java的程序肯定要选择一个java的base image）。然后我们根据不同类型的程序生成一个Dockerfile，把用户的代码和base image揉在一起，搞定！

# 安装方法

代码在 [这里](https://github.com/dinp/builder) ，这是个golang的项目，使用beego框架，安装起来也比较简单，就是通常的golang项目的编译方式

	mkdir -p $GOPATH/src/github.com/dinp # 假设你已经配置好了GOROOT和GOPATH
	cd $GOPATH/src/github.com/dinp
	git clone https://github.com/dinp/builder.git
	cd builder
	go get ./...
	go build
	mv conf/app.conf.example conf/app.conf
	# modify conf/app.conf
	./builder

项目采用[UIC](http://ulricqin.com/project/uic/)作为单点登录系统，所以要想运行测试的话得先把UIC搭建起来

简单解释一下`conf/app.conf`各项配置的作用

- appname、httpport没啥好说的
- runmode取值是dev或者prod，可以参看beego的文档
- db打头的是数据库配置，数据库初始化脚本在schema.sql
- tmpdir、logdir是就是一些临时目录，不解释
- buildtimeout是编译超时时间，超过了这个时间就会被kill，单位是分钟
- uicinternal、uicexternal是UIC的配置，为啥分成两个呢？Builder和UIC通常是在一个内网的，相互之间的访问可以走内网，所以有个uicinternal，但是有的时候UIC的内网地址用户是没法访问的，所以在sso登录的时候还是需要redirect到UIC的外网地址
- registry是docker私有源，读者可以使用docker registry搭建
- buildscript是Builder内部用到的一个脚本文件地址，默认配置即可
- tplmapping这个很重要，这是base image的配置，在registry中增加了base image，也要在此配置一下，这样用户才能在页面上看到
- token，这是与UIC通信的凭证，与UIC的token配置成一样即可


# 问题

编配的时候用户要查看log，这造成了单点状态，目前可以通过挂载分布式文件系统解决
