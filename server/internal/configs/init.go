package configs

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewShopAdminServerConfig,
	ParseWxPay,
	ParseJwt,
	ParseAuthnJwt,
)
