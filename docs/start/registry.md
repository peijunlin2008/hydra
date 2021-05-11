#  服务注册

每个`ServerType`有各自的服务注册函数，注册函数位于`hydra.S.xxxx`和`app`对象中，两种方式完全等价。

### 1. 静态注册

是指必须在服务器启动前完成的注册，即在`hydra.OnStarting`前执行。目前`API`,`Web`,`WS`,`RPC`只支持静态注册。

可以在以下阶段进行服务注册:

1. 包初始化```func init(){}```
2. app准备就绪```hydra.OnReadying(func() error {return nil})``` 或  ```hydra.OnReady(func()error {return nil})```
3. 服务准备启动前```hydra.S.OnSetup(func(conf app.IAPPConf) error {return nil})```

当注册发生错误时，可直接返回`error`,服务器将停止启动。

#### 1. 静态注册方法
|         ServerType         | app注册   | hydra.S注册   | 示例                                             |
| :------------------------: | :-------- | :------------ | :----------------------------------------------- |
|          http.API          | app.API   | hydra.S.API   | `app.API("/order/request",orders.OrderHandle)`   |
|          http.Web          | app.Web   | hydra.S.Web   | `app.Web("/order/request",orders.OrderHandle)`   |
|          http.WS           | app.WS    | hydra.S.WS    | `app.WS("/order/request",orders.OrderHandle)`    |
|          rpc.RPC           | app.RPC   | hydra.S.RPC   | `app.RPC("/order/request",orders.OrderHandle)`   |
|         cron.CRON          | app.CRON  | hydra.S.CRON  | `app.CRON("/order/request",orders.OrderHandle)`  |
|          mqc.MQC           | app.MQC   | hydra.S.MQC   | `app.MQC("/order/request",orders.OrderHandle)`   |
| http.API, http.Web,rpc.RPC | app.Micro | hydra.S.Micro | `app.Micro("/order/request",orders.OrderHandle)` |
|     cron.CRON,mqc.MQC      | app.Flow  | hydra.S.Flow  | `app.Flow("/order/request",orders.OrderHandle)`  |


服务注册时可指定与`ServerType`对应的参数

#### 2. 特殊参数
1. API服务编码方式

|函数|示例|
|:----:|:-----|
|app.API, app.Web,app.RPC,app.Micro |`app.Micro("/order/request",orders.OrderHandle,router.WithEncoding("gbk"))`|
|hydra.S.API, hydra.S.Web,hydra.S.RPC,hydra.S.Micro |`hydra.S.Micro("/order/request",orders.OrderHandle,router.WithEncoding("gbk"))`|


2. cron执行周期

|函数|示例|
|:----:|:-----|
|app.CRON|`app.CRON("/order/request",orders.OrderHandle,"@every 10s")`|
|hydra.S.CRON |`hydra.S.CRON("/order/request",orders.OrderHandle,"0 1 0 * *")`|


3. mqc队列名称

|函数|示例|
|:----:|:-----|
|app.MQC|`app.MQC("/order/request",orders.OrderHandle,"order:request")`|
|hydra.S.MQC |`hydra.S.MQC("/order/request",orders.OrderHandle,"order:query")`|


### 3. 动态注册
动态注册是指将`MQC`的服务注册与队列绑定，`CRON`的服务注册与周期表绑定分为两步，即服务注册与参数绑定进行分离

服务注册在服务器启动前执行，参数绑定在业务逻辑的任意阶段执行。


#### 1. MQC的动态注册

服务注册

```go
func init(){
    hydra.OnReady(func() {
        //注册服务 /order/request
        hydra.S.MQC("/order/request",orders.OrderHandle)
    })
}
```

参数绑定 
```go
hydra.MQC.Add("v1:order:request", "/order/request")

//同一服务绑定多个队列
hydra.MQC.Add("v2:order:request", "/order/request") 


//可指定并发执行协程数
hydra.MQC.Add("v3:order:request", "/order/request",100)
```

#### 2. CRON的动态注册

服务注册

```go
func init(){
    hydra.OnReady(func() {

        //注册服务 /order/request
        hydra.S.CRON("/order/request",orders.OrderHandle)
    })
}
```

参数绑定 
```go
hydra.CRON.Add("@now","/order/request") //立即执行1次

//同一服务可绑定多个执行周期
hydra.CRON.Add("@every 10s","/order/request") //每10秒执行一次


hydra.CRON.Add("0 12 0 * ?","/order/request") //每天12点执行1次

```


服务移除

```go

//根据对列名称和服务名移除
hydra.MQC.Remove("v2:order:request", "/order/request")

//根据定时表达式和服务名移除
hydra.CRON.Remove("0 12 0 * ?","/order/request") 

```