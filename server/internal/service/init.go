package service

import (
	"github.com/google/wire"
	"github.com/liujitcn/shop-admin/server/internal/service/admin"
	"github.com/liujitcn/shop-admin/server/internal/service/app"
	"github.com/liujitcn/shop-admin/server/internal/service/config"
	"github.com/liujitcn/shop-admin/server/internal/service/file"
	"github.com/liujitcn/shop-admin/server/internal/service/login"
	"github.com/liujitcn/shop-admin/server/internal/service/pay"
	payBiz "github.com/liujitcn/shop-admin/server/internal/service/pay/biz"

	adminBiz "github.com/liujitcn/shop-admin/server/internal/service/admin/biz"
	"github.com/liujitcn/shop-admin/server/internal/service/admin/task"
	appBiz "github.com/liujitcn/shop-admin/server/internal/service/app/biz"
	configBiz "github.com/liujitcn/shop-admin/server/internal/service/config/biz"
	fileBiz "github.com/liujitcn/shop-admin/server/internal/service/file/biz"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(

	adminBiz.NewBaseApiCase,
	adminBiz.NewBaseAreaCase,
	adminBiz.NewBaseConfigCase,
	adminBiz.NewBaseDeptCase,
	adminBiz.NewBaseDictCase,
	adminBiz.NewBaseDictItemCase,
	adminBiz.NewBaseJobCase,
	adminBiz.NewBaseJobLogCase,
	adminBiz.NewBaseLogCase,
	adminBiz.NewBaseMenuCase,
	adminBiz.NewBaseRoleCase,
	adminBiz.NewBaseUserCase,

	adminBiz.NewCasbinRuleCase,

	adminBiz.NewDashboardCase,

	adminBiz.NewGoodsCategoryCase,
	adminBiz.NewGoodsCase,
	adminBiz.NewGoodsPropCase,
	adminBiz.NewGoodsSkuCase,
	adminBiz.NewGoodsSpecCase,

	adminBiz.NewOrderCase,
	adminBiz.NewOrderAddressCase,
	adminBiz.NewOrderCancelCase,
	adminBiz.NewOrderGoodsCase,
	adminBiz.NewOrderLogisticsCase,
	adminBiz.NewOrderPaymentCase,
	adminBiz.NewOrderRefundCase,

	adminBiz.NewPayBillCase,

	adminBiz.NewShopBannerCase,
	adminBiz.NewShopHotCase,
	adminBiz.NewShopHotItemCase,
	adminBiz.NewShopServiceCase,

	adminBiz.NewUserStoreCase,

	appBiz.NewShopServiceCase,

	configBiz.NewConfigCase,
	fileBiz.NewFileCase,

	payBiz.NewOrderSchedulerCase,
	payBiz.NewPayCase,
	payBiz.NewPayBillCase,
	payBiz.NewWxPayCase,

	task.NewTradeBill,
	task.NewTaskList,

	admin.NewAuthService,
	admin.NewBaseApiService,
	admin.NewBaseConfigService,
	admin.NewBaseDeptService,
	admin.NewBaseDictService,
	admin.NewBaseJobService,
	admin.NewBaseLogService,
	admin.NewBaseMenuService,
	admin.NewBaseRoleService,
	admin.NewBaseUserService,

	admin.NewDashboardService,

	admin.NewGoodsCategoryService,
	admin.NewGoodsService,
	admin.NewGoodsPropService,
	admin.NewGoodsSkuService,
	admin.NewGoodsSpecService,

	admin.NewOrderService,

	admin.NewPayBillService,

	admin.NewShopBannerService,
	admin.NewShopHotService,
	admin.NewShopServiceService,

	admin.NewUserStoreService,

	app.NewShopServiceService,

	config.NewConfigService,

	file.NewFileService,

	login.NewLoginService,

	pay.NewPayService,
)
