package services

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
)

var suffixList = []string{defHandling, defHandler, defHandled, defFallback}

func reflectHandle(path string, h interface{}) (g *UnitGroup, err error) {
	//检查参数
	if path == "" || h == nil {
		return nil, fmt.Errorf("注册对象不能为空")
	}

	//输入参数为函数
	current := newUnitGroup(path)
	if vv, ok := h.(func(context.IContext) interface{}); ok {
		current.AddHandle("", context.Handler(vv))
		return current, nil
	}

	//检查输入的注册服务必须为struct
	typ := reflect.TypeOf(h)
	val := reflect.ValueOf(h)
	if val.Kind() == reflect.String {
		if _, ok := global.IsProto(h.(string), global.ProtoRPC); ok {
			current.AddHandle("", nil)
			return current, nil
		}
	}

	//检查传入的是构建函数
	if val.Kind() == reflect.Func {
		nval, err := createObject(h)
		if err != nil {
			return nil, err
		}
		typ = reflect.TypeOf(nval)
		val = reflect.ValueOf(nval)

	}

	//检查对象类型
	if val.Kind() != reflect.Ptr && val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("只能接收引用类型或struct; 实际是 %s", val.Kind())
	}

	//reflect所有函数，检查函数签名
	for i := 0; i < typ.NumMethod(); i++ {

		//检查函数参数是否符合接口要求
		mName := typ.Method(i).Name
		method := val.MethodByName(mName)

		//处理服务关闭函数
		if mName == defClose {
			current.Closing = method.Interface()
			continue
		}

		hasSuffix := checkSuffix(mName)

		//处理handling,handle,handled,fallback
		nfx, ok := method.Interface().(func(context.IContext) interface{})
		if !ok {
			if hasSuffix {
				err = fmt.Errorf("函数【%s】是钩子类型（%v）,但签名不是func(context.IContext) interface{}", mName, suffixList)
				return
			}
			continue
		}

		//检查函数名是否符合特定要求
		var nf context.Handler = nfx
		switch {
		case strings.HasSuffix(mName, defHandling):
			endName := strings.ToLower(mName[0 : len(mName)-len(defHandling)])
			current.AddHandling(endName, nf)
		case strings.HasSuffix(mName, defHandler):
			endName := strings.ToLower(mName[0 : len(mName)-len(defHandler)])
			current.AddHandle(endName, nf)
		case strings.HasSuffix(mName, defHandled):
			endName := strings.ToLower(mName[0 : len(mName)-len(defHandled)])
			current.AddHandled(endName, nf)
		case strings.HasSuffix(mName, defFallback):
			endName := strings.ToLower(mName[0 : len(mName)-len(defFallback)])
			current.AddFallback(endName, nf)
		}
	}
	if len(current.Services) == 0 {
		return nil, fmt.Errorf("%s中，未找到可用于注册的处理函数", path)
	}
	for _, u := range current.Services {
		if u.Handle == nil {
			return nil, fmt.Errorf("%s中,未指定[%s]的Handle函数", path, u.Service)
		}
	}
	return current, nil

}

func checkSuffix(mName string) bool {
	for i := range suffixList {
		if strings.HasSuffix(mName, suffixList[i]) {
			return true
		}
	}
	return false
}
