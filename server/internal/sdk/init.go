package sdk

import (
	"context"
	"time"

	"github.com/google/wire"
	authData "github.com/liujitcn/go-sdk/auth/data"
	"github.com/liujitcn/go-sdk/cache"
	"github.com/liujitcn/go-sdk/gorm"
	"github.com/liujitcn/go-sdk/oss"
	"github.com/liujitcn/go-sdk/queue"
	"github.com/liujitcn/shop-gorm-gen/models"
	authnEngine "github.com/tx7do/kratos-authn/engine"
	"github.com/tx7do/kratos-authn/engine/jwt"
	authzEngine "github.com/tx7do/kratos-authz/engine"
	authzEngineCasbin "github.com/tx7do/kratos-authz/engine/casbin"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	RegisterMigrateModels,
	queue.NewQueue,
	cache.NewCache,
	oss.NewOSS,
	gorm.NewGormClient,
	NewAuthenticator,
	NewAuthzEngine,
	NewUserToken,
)

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

// NewAuthenticator 创建认证器
func NewAuthenticator(cfg *conf.Authentication_Jwt) authnEngine.Authenticator {
	authenticator, _ := jwt.NewAuthenticator(
		jwt.WithKey([]byte(cfg.GetKey())),
		jwt.WithSigningMethod(cfg.GetMethod()),
	)
	return authenticator
}

// NewAuthzEngine 创建鉴权引擎
func NewAuthzEngine() (authzEngine.Engine, error) {
	return authzEngineCasbin.NewEngine(context.Background())
}

func NewUserToken(cache cache.Cache, authenticator authnEngine.Authenticator) *authData.UserToken {
	const (
		userAccessTokenKeyPrefix  = "uat_"
		userRefreshTokenKeyPrefix = "urt_"
	)
	return authData.NewUserToken(
		cache,
		authenticator,
		userAccessTokenKeyPrefix,
		userRefreshTokenKeyPrefix,
		2*time.Hour,
		24*7*time.Hour,
	)
}
