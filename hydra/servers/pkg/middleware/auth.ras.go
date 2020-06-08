package middleware

import (
	"fmt"
	"net/http"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/lib4go/types"
)

//RASAuth RAS远程认证
func RASAuth() Handler {
	return func(ctx IMiddleContext) {

		//获取FSA配置
		ras := ctx.ServerConf().GetRASConf()
		if ras.Disable {
			ctx.Next()
			return
		}

		b, auth := ras.Match(ctx.Request().Path().GetRouter().Service)
		if !b {
			ctx.Next()
			return
		}

		ctx.Response().AddSpecial("ras")

		input, err := ctx.Request().GetMap()
		if err != nil {
			ctx.Response().Abort(http.StatusInternalServerError, err)
			return
		}

		input["__auth_"], err = auth.AuthString()
		if err != nil {
			ctx.Response().Abort(http.StatusInternalServerError, err)
			return
		}

		respones, err := components.Def.RPC().GetRegularRPC().Request(ctx.Context(), auth.Service, input)
		if err != nil || !respones.Success() {
			ctx.Response().Abort(types.GetMax(respones.Status, http.StatusForbidden), fmt.Errorf("远程认证失败:%s,err:%v(%d)", err, respones.Result, respones.Status))
			return
		}
		result, err := respones.GetResult()
		if err != nil {
			ctx.Response().Abort(http.StatusForbidden, err)
			return
		}
		ctx.Meta().MergeMap(result)
		return
	}
}
