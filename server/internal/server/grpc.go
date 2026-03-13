package server

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	bootstrapConf "github.com/liujitcn/kratos-kit/api/gen/go/conf"
	"github.com/liujitcn/kratos-kit/auth"
	authnEngine "github.com/liujitcn/kratos-kit/auth/authn/engine"
	authzEngine "github.com/liujitcn/kratos-kit/auth/authz/engine"
	authData "github.com/liujitcn/kratos-kit/auth/data"
	"github.com/liujitcn/kratos-kit/bootstrap"
	"github.com/liujitcn/kratos-kit/rpc"
	adminApi "github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	configApi "github.com/liujitcn/shop-admin/server/api/gen/go/config"
	fileApi "github.com/liujitcn/shop-admin/server/api/gen/go/file"
	loginApi "github.com/liujitcn/shop-admin/server/api/gen/go/login"
	payApi "github.com/liujitcn/shop-admin/server/api/gen/go/pay"
	"github.com/liujitcn/shop-admin/server/internal/middleware/logging"
	"github.com/liujitcn/shop-admin/server/internal/service/admin"
	"github.com/liujitcn/shop-admin/server/internal/service/admin/biz"
	"github.com/liujitcn/shop-admin/server/internal/service/config"
	"github.com/liujitcn/shop-admin/server/internal/service/file"
	"github.com/liujitcn/shop-admin/server/internal/service/login"
	"github.com/liujitcn/shop-admin/server/internal/service/pay"
)

// GrpcMiddlewares 为 gRPC 服务注入专用中间件类型，避免与 HTTP 中间件冲突。
type GrpcMiddlewares []middleware.Middleware

// NewGrpcMiddleware 创建中间件
func NewGrpcMiddleware(
	ctx *bootstrap.Context,
	authenticator authnEngine.Authenticator,
	userCase *biz.BaseUserCase,
	authorizer authzEngine.Engine,
	userToken *authData.UserToken,
	jwtCfg *bootstrapConf.Authentication_Jwt,
) GrpcMiddlewares {
	var ms GrpcMiddlewares
	cfg := ctx.GetConfig()
	if cfg != nil && cfg.Server != nil && cfg.Server.Grpc != nil && cfg.Server.Grpc.Middleware != nil && cfg.Server.Grpc.Middleware.EnableLogging {
		ms = append(ms, logging.Server(ctx.GetLogger(), userCase, authenticator))
	}
	ms = append(ms, auth.NewAuthMiddleware(authenticator, authorizer, userToken, jwtCfg))
	return ms
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	ctx *bootstrap.Context,
	middlewares GrpcMiddlewares,

	adminAuth *admin.AuthService,
	adminBaseApi *admin.BaseApiService,
	adminBaseConfig *admin.BaseConfigService,
	adminBaseDept *admin.BaseDeptService,
	adminBaseDict *admin.BaseDictService,
	adminBaseJob *admin.BaseJobService,
	adminBaseLog *admin.BaseLogService,
	adminBaseMenu *admin.BaseMenuService,
	adminBaseRole *admin.BaseRoleService,
	adminBaseUser *admin.BaseUserService,

	adminDashboard *admin.DashboardService,

	adminGoodsCategory *admin.GoodsCategoryService,
	adminGoods *admin.GoodsService,
	adminGoodsProp *admin.GoodsPropService,
	adminGoodsSku *admin.GoodsSkuService,
	adminGoodsSpec *admin.GoodsSpecService,

	adminOrder *admin.OrderService,

	adminPayBill *admin.PayBillService,

	adminShopBanner *admin.ShopBannerService,
	adminShopHot *admin.ShopHotService,
	adminShopService *admin.ShopServiceService,

	adminUserStore *admin.UserStoreService,

	config *config.ConfigService,
	file *file.FileService,
	login *login.LoginService,
	pay *pay.PayService,
) (*grpc.Server, error) {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Http == nil {
		return nil, nil
	}

	srv, err := rpc.CreateGrpcServer(cfg, middlewares...)
	if err != nil {
		return nil, err
	}
	adminApi.RegisterAuthServiceServer(srv, adminAuth)
	adminApi.RegisterBaseApiServiceServer(srv, adminBaseApi)
	adminApi.RegisterBaseConfigServiceServer(srv, adminBaseConfig)
	adminApi.RegisterBaseDeptServiceServer(srv, adminBaseDept)
	adminApi.RegisterBaseDictServiceServer(srv, adminBaseDict)
	adminApi.RegisterBaseJobServiceServer(srv, adminBaseJob)
	adminApi.RegisterBaseLogServiceServer(srv, adminBaseLog)
	adminApi.RegisterBaseMenuServiceServer(srv, adminBaseMenu)
	adminApi.RegisterBaseRoleServiceServer(srv, adminBaseRole)
	adminApi.RegisterBaseUserServiceServer(srv, adminBaseUser)
	adminApi.RegisterDashboardServiceServer(srv, adminDashboard)
	adminApi.RegisterGoodsCategoryServiceServer(srv, adminGoodsCategory)
	adminApi.RegisterGoodsServiceServer(srv, adminGoods)
	adminApi.RegisterGoodsPropServiceServer(srv, adminGoodsProp)
	adminApi.RegisterGoodsSkuServiceServer(srv, adminGoodsSku)
	adminApi.RegisterGoodsSpecServiceServer(srv, adminGoodsSpec)
	adminApi.RegisterOrderServiceServer(srv, adminOrder)
	adminApi.RegisterPayBillServiceServer(srv, adminPayBill)
	adminApi.RegisterShopBannerServiceServer(srv, adminShopBanner)
	adminApi.RegisterShopHotServiceServer(srv, adminShopHot)
	adminApi.RegisterShopServiceServiceServer(srv, adminShopService)
	adminApi.RegisterUserStoreServiceServer(srv, adminUserStore)
	configApi.RegisterConfigServiceServer(srv, config)
	fileApi.RegisterFileServiceServer(srv, file)
	loginApi.RegisterLoginServiceServer(srv, login)
	payApi.RegisterPayServiceServer(srv, pay)

	return srv, nil
}
