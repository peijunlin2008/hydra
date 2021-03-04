package ws

import (
	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
)

var wsInternalEngine *wsEngine
var metric = middleware.NewMetric()

//InitWSEngine 创建默认的WS引擎
func InitWSEngine(routers ...*router.Router) {
	wsInternalEngine = newWSEngine(routers...)
}

type wsEngine struct {
	*dispatcher.Engine
	metric        *middleware.Metric
	adapterEngine *adapter.Engine
}

func newWSEngine(routers ...*router.Router) *wsEngine {
	s := &wsEngine{
		metric: metric,
	}
	s.adapterEngine = adapter.New()
	s.Engine = s.adapterEngine.DispEngine()

	s.adapterEngine.Use(middleware.Recovery())
	s.adapterEngine.Use(middleware.Logging()) //记录请求日志
	s.adapterEngine.Use(middleware.Recovery())
	s.adapterEngine.Use(middleware.Tag())
	s.adapterEngine.Use(middleware.Trace()) //跟踪信息
	s.adapterEngine.Use(middleware.Limit()) //限流处理
	s.adapterEngine.Use(middleware.Delay()) //
	s.adapterEngine.Use(middleware.APIKeyAuth())
	s.adapterEngine.Use(middleware.RASAuth())
	s.adapterEngine.Use(middleware.JwtAuth())   //jwt安全认证
	s.adapterEngine.Use(middleware.Render())    //响应渲染组件
	s.adapterEngine.Use(middleware.JwtWriter()) //设置jwt回写
	s.adapterEngine.Use(middlewares...)
	s.adapterEngine.Use(s.metric.Handle()) //生成metric报表

	s.addWSRouter(routers...)
	return s
}
func (s *wsEngine) addWSRouter(routers ...*router.Router) {
	adapterRouters := make([]adapter.IRouter, len(routers))
	for i := range routers {
		adapterRouters[i] = routers[i]
	}
	s.adapterEngine.DispHandle(global.WS, adapterRouters...)
}
