package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liujitcn/shop-admin/server/api/gen/go/conf"
	"github.com/liujitcn/shop-admin/server/internal/configs"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	bootstrapConf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	//_ "github.com/tx7do/kratos-bootstrap/config/apollo"
	//_ "github.com/tx7do/kratos-bootstrap/config/consul"
	//_ "github.com/tx7do/kratos-bootstrap/config/etcd"
	//_ "github.com/tx7do/kratos-bootstrap/config/kubernetes"
	//_ "github.com/tx7do/kratos-bootstrap/config/nacos"
	//_ "github.com/tx7do/kratos-bootstrap/config/polaris"

	//_ "github.com/tx7do/kratos-bootstrap/logger/aliyun"
	//_ "github.com/tx7do/kratos-bootstrap/logger/fluent"
	//_ "github.com/tx7do/kratos-bootstrap/logger/logrus"
	//_ "github.com/tx7do/kratos-bootstrap/logger/tencent"
	_ "github.com/tx7do/kratos-bootstrap/logger/zap"
	//_ "github.com/tx7do/kratos-bootstrap/logger/zerolog"
	//_ "github.com/tx7do/kratos-bootstrap/registry/consul"
	//_ "github.com/tx7do/kratos-bootstrap/registry/etcd"
	//_ "github.com/tx7do/kratos-bootstrap/registry/eureka"
	//_ "github.com/tx7do/kratos-bootstrap/registry/kubernetes"
	//_ "github.com/tx7do/kratos-bootstrap/registry/nacos"
	//_ "github.com/tx7do/kratos-bootstrap/registry/polaris"
	//_ "github.com/tx7do/kratos-bootstrap/registry/servicecomb"
	//_ "github.com/tx7do/kratos-bootstrap/registry/zookeeper"
)

var (
	Project = "shop"
	AppId   = "admin"
	version = "1.0.0"
)

func newApp(
	ctx *bootstrap.Context,
	hs *http.Server,
) *kratos.App {
	return bootstrap.NewApp(ctx,
		hs,
	)
}

func runApp() error {
	ctx := bootstrap.NewContext(
		context.Background(),
		&bootstrapConf.AppInfo{
			Project: Project,
			AppId:   AppId,
			Version: version,
		},
	)
	ctx.RegisterCustomConfig(configs.WrapperConfigKey, &conf.ShopAdminServerConfigWrapper{})
	return bootstrap.RunApp(ctx, initApp)
}

func main() {
	if err := runApp(); err != nil {
		panic(err)
	}
}
