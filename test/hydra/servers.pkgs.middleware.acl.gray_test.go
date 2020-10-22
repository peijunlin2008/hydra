package hydra

import (
	"fmt"
	xhttp "net/http"
	"os"
	"testing"
	"time"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

//author:liujinyin
//time:2020-10-20 15:30
//desc:测试灰度中间件逻辑
func TestGray_Disable(t *testing.T) {

	type testCase struct {
		name        string
		opts        string
		upservers   []string
		serviceAddr string
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{
			name:        "灰度-未启用-未配置",
			opts:        ``,
			serviceAddr: "http://localhost:51002",
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
		{
			name:        "灰度-未启用-配置为关闭",
			opts:        `{"disable":true}`,
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
	}
	for _, tt := range tests {
		global.Def.ServerTypes = []string{http.API}
		fmt.Println("---------------------------", tt.name)

		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		//mockConf.Service.API.Add()
		//初始化测试用例参数
		mockConf.GetAPI().Gray(tt.opts)
		serverConf := mockConf.GetAPIConf()

		request, _ := xhttp.NewRequest(xhttp.MethodGet, "http://localhost:51001/upcluster", nil)
		request.Header = xhttp.Header{}
		request.Header.Add("content-type", "text/plain")

		ctx := &mocks.MiddleContext{
			MockTFuncs:  map[string]interface{}{},
			HttpRequest: request,
			HttpResponse: &mocks.MockResponseWriter{
				ResponseHeader: xhttp.Header{},
			},
			MockResponse:   &mocks.MockResponse{MockStatus: 200},
			MockServerConf: serverConf,
		}

		//获取中间件
		handler := middleware.Gray()

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()

		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantContent, gotContent)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)

	}
}

//author:liujinyin
//time:2020-10-20 15:30
//desc:测试灰度中间件逻辑
func TestGray_Enable_Has(t *testing.T) {
	startUpstreamServer(":51003")

	type testCase struct {
		name        string
		opts        string
		upservers   []string
		serviceAddr string
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{
			name:        "灰度-启用-模板匹配=>false",
			opts:        `{"disable":false,"filter":"false","upcluster":"t"}`,
			wantStatus:  200,
			wantContent: "",
			wantSpecial: "",
		},
		{
			name:        "灰度-启用-模板匹配需要-有地址可用",
			opts:        `{"disable":false,"filter":"true","upcluster":"t"}`,
			wantStatus:  305,
			wantContent: "",
			wantSpecial: "gray",
		},
	}
	for _, tt := range tests {
		global.Def.ServerTypes = []string{http.API}
		fmt.Println("---------------------------", tt.name)

		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		//mockConf.Service.API.Add()
		//初始化测试用例参数
		mockConf.GetAPI().Gray(tt.opts)
		serverConf := mockConf.GetAPIConf()

		request, _ := xhttp.NewRequest(xhttp.MethodGet, "http://localhost:51001/upcluster", nil)
		request.Header = xhttp.Header{}
		request.Header.Add("content-type", "text/plain")

		ctx := &mocks.MiddleContext{
			MockTFuncs:  map[string]interface{}{},
			HttpRequest: request,
			HttpResponse: &mocks.MockResponseWriter{
				ResponseHeader: xhttp.Header{},
			},
			MockResponse:   &mocks.MockResponse{MockStatus: 200},
			MockServerConf: serverConf,
		}

		//获取中间件
		handler := middleware.Gray()

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()

		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantContent, gotContent)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)

	}
}

//author:liujinyin
//time:2020-10-20 15:30
//desc:测试灰度中间件逻辑
func TestGray_Enable_None(t *testing.T) {
	type testCase struct {
		name        string
		opts        string
		upservers   []string
		serviceAddr string
		wantStatus  int
		wantContent string
		wantSpecial string
	}

	tests := []*testCase{
		{
			name:        "灰度-启用-模板匹配需要-无上游地址",
			opts:        `{"disable":false,"filter":"true","upcluster":"t"}`,
			wantStatus:  502,
			wantContent: "",
			wantSpecial: "gray",
		},
		// {
		// 	name:        "灰度-启用-模板匹配需要-有地址不可用",
		// 	opts:        `{"disable":true,"filter":"true","upcluster":"t"}`,
		// 	wantStatus:  200,
		// 	wantContent: "",
		// 	wantSpecial: "",
		// },
	}
	for _, tt := range tests {
		global.Def.ServerTypes = []string{http.API}
		fmt.Println("---------------------------", tt.name)

		mockConf := mocks.NewConf()
		mockConf.API(":51001")
		//mockConf.Service.API.Add()
		//初始化测试用例参数
		mockConf.GetAPI().Gray(tt.opts)
		serverConf := mockConf.GetAPIConf()

		request, _ := xhttp.NewRequest(xhttp.MethodGet, "http://localhost:51001/upcluster", nil)
		request.Header = xhttp.Header{}
		request.Header.Add("content-type", "text/plain")

		ctx := &mocks.MiddleContext{
			MockTFuncs:  map[string]interface{}{},
			HttpRequest: request,
			HttpResponse: &mocks.MockResponseWriter{
				ResponseHeader: xhttp.Header{},
			},
			MockResponse:   &mocks.MockResponse{MockStatus: 200},
			MockServerConf: serverConf,
		}

		//获取中间件
		handler := middleware.Gray()

		//调用中间件
		handler(ctx)

		//断言结果
		gotStatus, gotContent := ctx.Response().GetFinalResponse()
		gotSpecial := ctx.Response().GetSpecials()

		assert.Equalf(t, tt.wantStatus, gotStatus, tt.name, tt.wantStatus, gotStatus)
		assert.Equalf(t, tt.wantContent, gotContent, tt.name, tt.wantContent, gotContent)
		assert.Equalf(t, tt.wantSpecial, gotSpecial, tt.name, tt.wantSpecial, gotSpecial)

	}
}

func startUpstreamServer(port string) {
	app := hydra.NewApp(
		hydra.WithPlatName("hydra"),
		hydra.WithSystemName("apiserver"),
		hydra.WithServerTypes(http.API),
		hydra.WithClusterName("t"),
		//hydra.WithRegistry("zk://192.168.0.101"),
		hydra.WithRegistry("lm://."),
	)
	hydra.Conf.API(port)
	app.API("/upcluster", upcluster)

	os.Args = []string{"upclusterserver", "run"}
	go app.Start()
	time.Sleep(time.Second * 10)
}

func upcluster(ctx hydra.IContext) interface{} {
	return "upcluster"
}
