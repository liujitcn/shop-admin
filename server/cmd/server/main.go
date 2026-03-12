package main

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liujitcn/kratos-kit/bootstrap"
	"github.com/liujitcn/shop-admin/server/api/gen/go/conf"
	"github.com/liujitcn/shop-admin/server/internal/configs"
	"github.com/liujitcn/shop-gorm-gen/models"

	bootstrapConf "github.com/liujitcn/kratos-kit/api/gen/go/conf"

	//_ "github.com/liujitcn/kratos-kit/database/gorm/driver/bigquery"
	_ "github.com/liujitcn/kratos-kit/database/gorm/driver/mysql"
	//_ "github.com/liujitcn/kratos-kit/database/gorm/driver/oracle"
	//_ "github.com/liujitcn/kratos-kit/database/gorm/driver/postgres"
	//_ "github.com/liujitcn/kratos-kit/database/gorm/driver/sqlite"
	//_ "github.com/liujitcn/kratos-kit/database/gorm/driver/sqlserver"

	//_ "github.com/liujitcn/kratos-kit/config/apollo"
	//_ "github.com/liujitcn/kratos-kit/config/consul"
	//_ "github.com/liujitcn/kratos-kit/config/etcd"
	//_ "github.com/liujitcn/kratos-kit/config/kubernetes"
	//_ "github.com/liujitcn/kratos-kit/config/nacos"
	//_ "github.com/liujitcn/kratos-kit/config/polaris"

	//_ "github.com/liujitcn/kratos-kit/logger/aliyun"
	//_ "github.com/liujitcn/kratos-kit/logger/fluent"
	//_ "github.com/liujitcn/kratos-kit/logger/logrus"
	//_ "github.com/liujitcn/kratos-kit/logger/tencent"
	_ "github.com/liujitcn/kratos-kit/logger/zap"
	//_ "github.com/liujitcn/kratos-kit/logger/zerolog"
	//_ "github.com/liujitcn/kratos-kit/registry/consul"
	//_ "github.com/liujitcn/kratos-kit/registry/etcd"
	//_ "github.com/liujitcn/kratos-kit/registry/eureka"
	//_ "github.com/liujitcn/kratos-kit/registry/kubernetes"
	//_ "github.com/liujitcn/kratos-kit/registry/nacos"
	//_ "github.com/liujitcn/kratos-kit/registry/polaris"
	//_ "github.com/liujitcn/kratos-kit/registry/servicecomb"
	//_ "github.com/liujitcn/kratos-kit/registry/zookeeper"
)

var (
	Project = "shop"
	AppId   = "admin"
	version = "1.0.0"
)

func newApp(
	ctx *bootstrap.Context,
	gs *grpc.Server,
	hs *http.Server,
) *kratos.App {
	return bootstrap.NewApp(ctx,
		gs,
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

// RegisterMigrateModels registers all GORM models for migration.
func RegisterMigrateModels() []interface{} {
	return []interface{}{
		new(models.BaseAPI),
		new(models.BaseArea),
		new(models.BaseConfig),
		new(models.BaseDept),
		new(models.BaseDict),
		new(models.BaseDictItem),
		new(models.BaseJob),
		new(models.BaseJobLog),
		new(models.BaseLog),
		new(models.BaseMenu),
		new(models.BaseRole),
		new(models.BaseUser),
		new(models.CasbinRule),
		new(models.Goods),
		new(models.GoodsCategory),
		new(models.GoodsProp),
		new(models.GoodsSku),
		new(models.GoodsSpec),
		new(models.Order),
		new(models.OrderAddress),
		new(models.OrderCancel),
		new(models.OrderGoods),
		new(models.OrderLogistics),
		new(models.OrderPayment),
		new(models.OrderRefund),
		new(models.PayBill),
		new(models.ShopBanner),
		new(models.ShopHot),
		new(models.ShopHotGoods),
		new(models.ShopHotItem),
		new(models.ShopService),
		new(models.UserAddress),
		new(models.UserCart),
		new(models.UserCollect),
		new(models.UserStore),
	}
}
