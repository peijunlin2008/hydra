
# 参数配置

指服务启动、运行所必须的配置。

这些配置使用注册中心集中管理，本地零配置。

配置发生变化后自动通知到集群并立即生效，必要时会自动重启。

有三种方式构建配置：


### 代码构建
可通过`hydra.Conf.[ServerType]`提供的函数对服务参数与服务组件进行构建，

代码构建的配置必须通过配置安装命令`conf install -r ...`，安装到注册中心方可使用(当为单机模式，

则无须安装配置。单机模式即注册中心为本地内存，即`-r lm://.`,启动`run`,`start`,`install`时未指定`-r`时，系统默认按`-r lm://.`运行）




代码构建配置：
```go
hydra.Conf.API("6689")
```

部分配置构建时可能会用到平台名称、系统名等初始参数，建议将配置的初始化放到系统准备就绪后执行:

```go
func init() {
	hydra.OnReady(func() {
        hydra.Conf.API("6689")
    }
}
```
或

```go
func init() {
	hydra.OnReadying(func() {
        hydra.Conf.API("6689")
    }
}
```

> 注意：两种方式都是在系统参数检查完毕后立即运行，不同的是`OnReadying`先运行，`OnReady`后运行。

> `OnReadying` 和 `OnReady` 都属于系统勾子函数，了解其它勾子函数,请参数`勾子函数`章节

#### 1. 支持链式调用与交互参数输入(hydra.ByInstall):

```go
//go:embed loginweb/dist/static
var staticFs embed.FS
var archive = "loginweb/dist/static"

hydra.OnReady(func() {
    hydra.Conf.Web("6687", api.WithTimeout(300, 300), api.WithDNS(hydra.ByInstall)).
			Static(static.WithAutoRewrite(), static.WithEmbed(archive, staticFs)).
			Processor(processor.WithServicePrefix("/web")).
			Header(header.WithCrossDomain()).
			Jwt(jwt.WithMode("HS512"),
				jwt.WithSecret("f0abd74b09bcc61449d66ae5d8128c18"),
				jwt.WithExpireAt(36000),
				jwt.WithAuthURL("/"),
				jwt.WithHeader(),
				jwt.WithExcludes(
					"/*/system/config/get",
					"/*/member/login",
					"/*/member/bind/*",
					"/*/member/sendcode",
					"/*/logout",
				),
			)
})
```

> hydra.ByInstall 会在配置安装时通过向导方式引导用户输入

> 其它`ServerType`可通过`hydra.Conf.API`,`hydra.Conf.CRON`,`hydra.Conf.MQC`等进行参数配置

####  2. 公共参数配置
指同一个平台下所有`ServerType`可共同使用的配置, 通过```hydra.Conf.Vars()....```指定


DB配置：

```go
hydra.OnReady(func() {
    hydra.Conf.Vars().DB().MySQLByConnStr("db", "hydra:123456@tcp(192.168.0.36:10036)/hydra?charset=utf8")
})
```
缓存配置：

```go
hydra.OnReady(func() {
    hydra.Conf.Vars().Cache().GoCache("cache")
})
```

用户自定义配置:

```go
hydra.OnReady(func() {
    hydra.Conf.Vars().Custom("app", "conf", &model.LoginConf{
        LoginHost: hydra.ByInstall,
        APIHost:   hydra.ByInstall,
        UserLoginFailLimit: 5,
        UserLockTime:       24 * 60 * 60,
    })
})
```



### 文件构建



### 手动配置


