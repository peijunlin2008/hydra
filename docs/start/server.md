# 服务类型

目前支持六种服务类型(`ServerType`),多个`ServerType`可集成到同一个`app`。

 - 支持的`ServerType`:


| ServerType | 引用方式  | 说明                                         |
| :--------: | --------- | -------------------------------------------- |
|    API     | http.API  | 提供http服务                                 |
|    Web     | http.Web  | 提供Http服务                                 |
|     WS     | http.WS   | 提供websocket服务                            |
|    RPC     | rpc.RPC   | 提供基于grpc协议的通用远程调用服务           |
|    CRON    | cron.CRON | 提供定时任务服务，指定cron表达式定时执行任务 |
|    MQC     | mqc.MQC   | 提供消息消费服务，即message queue consumer   |


- 构建`app`:
```go
app := hydra.NewApp()
```

- 指定`ServerType`:

```go
app := hydra.NewApp(
        hydra.WithPlatName("test"),
        hydra.WithServerTypes(http.API,http.Web,cron,CRON,mqc.MQC,rpc.RPC,http.WS),
 )

```

- 完整示例

```go
package main

import (
    "github.com/micro-plat/hydra"   
    "github.com/micro-plat/hydra/hydra/servers/http"
    "github.com/micro-plat/hydra/hydra/servers/cron"
)

func main() {
	app := hydra.NewApp(
        hydra.WithPlatName("test"),
        hydra.WithServerTypes(http.API,cron.CRON),
    )
    app.API("/hello",hello)
    app.CRON("/auto",hello,"@every 5s") 
    app.Start()
}
func hello(ctx hydra.IContext) interface{} {
    return "hello world"
}
```
