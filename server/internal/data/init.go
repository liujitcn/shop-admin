package data

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	//genData.NewData,
	//genData.NewTransaction,
	NewBaseApiRepo,
	NewBaseAreaRepo,
	NewBaseConfigRepo,
	NewBaseDeptRepo,
	NewBaseDictRepo,
	NewBaseDictItemRepo,
	NewBaseJobRepo,
	NewBaseJobLogRepo,
	NewBaseLogRepo,
	NewBaseMenuRepo,
	NewBaseRoleRepo,
	NewBaseUserRepo,

	NewCasbinRuleRepo,

	NewGoodsCategoryRepo,
	NewGoodsRepo,
	NewGoodsPropRepo,
	NewGoodsSpecRepo,
	NewGoodsSkuRepo,

	NewOrderRepo,
	NewOrderAddressRepo,
	NewOrderCancelRepo,
	NewOrderGoodsRepo,
	NewOrderLogisticsRepo,
	NewOrderPaymentRepo,
	NewOrderRefundRepo,

	NewPayBillRepo,

	NewShopBannerRepo,
	NewShopHotRepo,
	NewShopHotGoodsRepo,
	NewShopHotItemRepo,
	NewShopServiceRepo,

	NewUserStoreRepo,
)
