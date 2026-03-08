package server

import (
	"io/fs"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liujitcn/go-sdk/auth"
	authData "github.com/liujitcn/go-sdk/auth/data"
	adminApi "github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	appApi "github.com/liujitcn/shop-admin/server/api/gen/go/app"
	"github.com/liujitcn/shop-admin/server/api/gen/go/conf"
	configApi "github.com/liujitcn/shop-admin/server/api/gen/go/config"
	loginApi "github.com/liujitcn/shop-admin/server/api/gen/go/login"
	payApi "github.com/liujitcn/shop-admin/server/api/gen/go/pay"
	"github.com/liujitcn/shop-admin/server/cmd/server/assets"
	customLogging "github.com/liujitcn/shop-admin/server/internal/middleware/logging"
	"github.com/liujitcn/shop-admin/server/internal/service"
	"github.com/liujitcn/shop-admin/server/internal/service/admin"
	"github.com/liujitcn/shop-admin/server/internal/service/admin/biz"
	"github.com/liujitcn/shop-admin/server/internal/service/app"
	"github.com/liujitcn/shop-admin/server/internal/service/config"
	"github.com/liujitcn/shop-admin/server/internal/service/file"
	"github.com/liujitcn/shop-admin/server/internal/service/login"
	"github.com/liujitcn/shop-admin/server/internal/service/pay"
	authnEngine "github.com/tx7do/kratos-authn/engine"
	authzEngine "github.com/tx7do/kratos-authz/engine"

	swaggerUI "github.com/tx7do/kratos-swagger-ui"

	"github.com/tx7do/kratos-bootstrap/rpc"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
)

// NewHttpMiddleware 创建中间件
func NewHttpMiddleware(
	ctx *bootstrap.Context,
	authenticator authnEngine.Authenticator,
	userCase *biz.BaseUserCase,
	authorizer authzEngine.Engine,
	userToken *authData.UserToken,
	jwtCfg *conf.Jwt,
) []middleware.Middleware {
	whiteList := jwtCfg.GetWhiteList()

	var ms []middleware.Middleware
	ms = append(ms, customLogging.Server(ctx.GetLogger(), userCase, authenticator))
	ms = append(ms, auth.NewAuthMiddleware(authenticator, authorizer, userToken, &auth.WhiteList{
		Prefix: whiteList.GetPrefix(),
		Regex:  whiteList.GetRegex(),
		Path:   whiteList.GetPath(),
		Match:  whiteList.GetMatch(),
	}))
	return ms
}

// NewHttpServer new an Http server.
func NewHttpServer(
	ctx *bootstrap.Context,
	middlewares []middleware.Middleware,
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

	appShopService *app.ShopServiceService,

	config *config.ConfigService,
	file *file.FileService,
	login *login.LoginService,
	pay *pay.PayService,
) (*http.Server, error) {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil, nil
	}

	srv, err := rpc.CreateRestServer(cfg,
		middlewares...,
	)
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

	appApi.RegisterShopServiceServiceHTTPServer(srv, appShopService)

	configApi.RegisterConfigServiceHTTPServer(srv, config)
	// 修改http接口实现
	service.RegisterFileServiceHTTPServer(srv, file)
	loginApi.RegisterLoginServiceHTTPServer(srv, login)
	payApi.RegisterPayServiceHTTPServer(srv, pay)

	if webFS, subErr := fs.Sub(assets.WebAssets, "web"); subErr == nil {
		webHandler := stdhttp.FileServer(stdhttp.FS(webFS))
		srv.HandlePrefix("/web/", stdhttp.StripPrefix("/web/", webHandler))
		srv.HandleFunc("/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			if r.URL.Path == "/" {
				stdhttp.Redirect(w, r, "/web/", stdhttp.StatusFound)
				return
			}
			stdhttp.NotFound(w, r)
		})
	}

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle(ctx.GetAppInfo().GetName()),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv, nil
}
