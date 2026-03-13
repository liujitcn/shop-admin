package server

import (
	"io/fs"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liujitcn/kratos-kit/auth"
	authnEngine "github.com/liujitcn/kratos-kit/auth/authn/engine"
	authzEngine "github.com/liujitcn/kratos-kit/auth/authz/engine"
	authData "github.com/liujitcn/kratos-kit/auth/data"
	adminApi "github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	configApi "github.com/liujitcn/shop-admin/server/api/gen/go/config"
	loginApi "github.com/liujitcn/shop-admin/server/api/gen/go/login"
	payApi "github.com/liujitcn/shop-admin/server/api/gen/go/pay"
	"github.com/liujitcn/shop-admin/server/cmd/server/assets"
	"github.com/liujitcn/shop-admin/server/internal/middleware/logging"
	"github.com/liujitcn/shop-admin/server/internal/service"
	"github.com/liujitcn/shop-admin/server/internal/service/admin"
	"github.com/liujitcn/shop-admin/server/internal/service/admin/biz"
	"github.com/liujitcn/shop-admin/server/internal/service/config"
	"github.com/liujitcn/shop-admin/server/internal/service/file"
	"github.com/liujitcn/shop-admin/server/internal/service/login"
	"github.com/liujitcn/shop-admin/server/internal/service/pay"

	bootstrapConf "github.com/liujitcn/kratos-kit/api/gen/go/conf"
	swaggerUI "github.com/liujitcn/kratos-kit/swagger-ui"

	"github.com/liujitcn/kratos-kit/rpc"

	"github.com/liujitcn/kratos-kit/bootstrap"
)

// HttpMiddlewares 为 HTTP 服务注入专用中间件类型，避免与 gRPC 中间件冲突。
type HttpMiddlewares []middleware.Middleware

// NewHttpMiddleware 创建中间件
func NewHttpMiddleware(
	ctx *bootstrap.Context,
	authenticator authnEngine.Authenticator,
	userCase *biz.BaseUserCase,
	authorizer authzEngine.Engine,
	userToken *authData.UserToken,
	jwtCfg *bootstrapConf.Authentication_Jwt,
) HttpMiddlewares {
	var ms HttpMiddlewares
	cfg := ctx.GetConfig()
	if cfg != nil && cfg.Server != nil && cfg.Server.Http != nil && cfg.Server.Http.Middleware != nil && cfg.Server.Http.Middleware.EnableLogging {
		ms = append(ms, logging.Server(ctx.GetLogger(), userCase, authenticator))
	}
	ms = append(ms, auth.NewAuthMiddleware(authenticator, authorizer, userToken, jwtCfg))
	return ms
}

// NewHttpServer new an Http server.
func NewHttpServer(
	ctx *bootstrap.Context,
	middlewares HttpMiddlewares,
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
) (*http.Server, error) {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Http == nil {
		return nil, nil
	}

	srv, err := rpc.CreateHttpServer(cfg, middlewares...)
	if err != nil {
		return nil, err
	}

	adminApi.RegisterAuthServiceHTTPServer(srv, adminAuth)
	adminApi.RegisterBaseApiServiceHTTPServer(srv, adminBaseApi)
	adminApi.RegisterBaseConfigServiceHTTPServer(srv, adminBaseConfig)
	adminApi.RegisterBaseDeptServiceHTTPServer(srv, adminBaseDept)
	adminApi.RegisterBaseDictServiceHTTPServer(srv, adminBaseDict)
	adminApi.RegisterBaseJobServiceHTTPServer(srv, adminBaseJob)
	adminApi.RegisterBaseLogServiceHTTPServer(srv, adminBaseLog)
	adminApi.RegisterBaseMenuServiceHTTPServer(srv, adminBaseMenu)
	adminApi.RegisterBaseRoleServiceHTTPServer(srv, adminBaseRole)
	adminApi.RegisterBaseUserServiceHTTPServer(srv, adminBaseUser)
	adminApi.RegisterDashboardServiceHTTPServer(srv, adminDashboard)
	adminApi.RegisterGoodsCategoryServiceHTTPServer(srv, adminGoodsCategory)
	adminApi.RegisterGoodsServiceHTTPServer(srv, adminGoods)
	adminApi.RegisterGoodsPropServiceHTTPServer(srv, adminGoodsProp)
	adminApi.RegisterGoodsSkuServiceHTTPServer(srv, adminGoodsSku)
	adminApi.RegisterGoodsSpecServiceHTTPServer(srv, adminGoodsSpec)
	adminApi.RegisterOrderServiceHTTPServer(srv, adminOrder)
	adminApi.RegisterPayBillServiceHTTPServer(srv, adminPayBill)
	adminApi.RegisterShopBannerServiceHTTPServer(srv, adminShopBanner)
	adminApi.RegisterShopHotServiceHTTPServer(srv, adminShopHot)
	adminApi.RegisterShopServiceServiceHTTPServer(srv, adminShopService)
	adminApi.RegisterUserStoreServiceHTTPServer(srv, adminUserStore)
	configApi.RegisterConfigServiceHTTPServer(srv, config)
	// 修改http接口实现
	service.RegisterFileServiceHTTPServer(srv, file)
	loginApi.RegisterLoginServiceHTTPServer(srv, login)
	payApi.RegisterPayServiceHTTPServer(srv, pay)

	if webFS, subErr := fs.Sub(assets.WebAssets, "web"); subErr == nil {
		webHandler := stdhttp.FileServer(stdhttp.FS(webFS))
		srv.HandlePrefix("/web/", stdhttp.StripPrefix("/web/", webHandler))
		srv.HandlePrefix("/js/", webHandler)
		srv.HandlePrefix("/css/", webHandler)
		srv.HandlePrefix("/img/", webHandler)
		srv.Handle("/favicon.ico", webHandler)
		srv.HandleFunc("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			if r.URL.Path == "/" {
				stdhttp.ServeFileFS(w, r, webFS, "index.html")
				return
			}
			stdhttp.NotFound(w, r)
		})
	}

	if cfg.GetServer().GetHttp().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle(ctx.GetAppInfo().GetName()),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv, nil
}
