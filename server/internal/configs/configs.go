package configs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/liujitcn/go-sdk"
	"github.com/liujitcn/shop-admin/server/api/gen/go/conf"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	bootstrapConf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/durationpb"
)

const WrapperConfigKey = "ShopAdminServer"

func NewShopAdminServerConfig(ctx *bootstrap.Context) *conf.ShopAdminServerConfig {
	cfg, ok := ctx.GetCustomConfig(WrapperConfigKey)
	if ok {
		wrapperCfg := cfg.(*conf.ShopAdminServerConfigWrapper)
		return wrapperCfg.GetShopAdminServer()
	}
	return &conf.ShopAdminServerConfig{}
}

func ParseAuthnJwt(ctx *bootstrap.Context) *bootstrapConf.Authentication_Jwt {
	cfg := ctx.GetConfig()
	if cfg == nil || cfg.GetAuthn() == nil || cfg.GetAuthn().GetJwt() == nil {
		return &bootstrapConf.Authentication_Jwt{
			Method: "HS256",
			Key:    "shop-admin",
		}
	}
	return cfg.GetAuthn().GetJwt()
}

func ParseJwt(cfg *conf.ShopAdminServerConfig) *conf.Jwt {
	jwt := cfg.GetJwt()
	if jwt == nil {
		jwt = &conf.Jwt{
			AccessTokenExpires:  durationpb.New(7200 * time.Second),
			RefreshTokenExpires: durationpb.New(14400 * time.Second),
			WhiteList:           &conf.WhiteList{},
		}
	}
	return jwt
}

func ParseWxPay(cfg *conf.ShopAdminServerConfig) (*conf.WxPay, error) {
	wxPay := cfg.GetWxPay()
	if wxPay == nil {
		return nil, errors.New("支付配置信息错误")
	}
	appid := wxPay.GetAppid()
	mchId := wxPay.GetMchId()
	mchCertSn := wxPay.GetMchCertSn()
	mchCertPath := wxPay.GetMchCertPath()
	mchAPIv3Key := wxPay.GetMchAPIv3Key()
	if appid == "" || mchId == "" || mchCertSn == "" || mchCertPath == "" || mchAPIv3Key == "" {
		return nil, errors.New("支付配置信息错误")
	}
	// 兼容不同工作目录启动（GoLand/命令行）导致的相对路径差异。
	if resolvedPath, ok := resolveFilePath(mchCertPath); ok {
		wxPay.MchCertPath = resolvedPath
	}
	return wxPay, nil
}

func resolveFilePath(path string) (string, bool) {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
		return path, false
	}

	candidates := []string{
		path,
		filepath.Join("server", path),
		filepath.Join("..", path),
		filepath.Join("..", "..", path),
		filepath.Join("..", "..", "..", path),
		filepath.Join("..", "server", path),
		filepath.Join(filepath.Dir(os.Args[0]), "..", path),
		filepath.Join(filepath.Dir(os.Args[0]), "..", "..", path),
	}

	for _, p := range candidates {
		cleaned := filepath.Clean(p)
		if _, err := os.Stat(cleaned); err == nil {
			return cleaned, true
		}
	}
	return path, false
}

func ParsePayTimeout() time.Duration {
	cache := sdk.Runtime.GetCache()
	if cache == nil {
		// 默认30分钟
		return time.Duration(_const.PayTimeout) * time.Minute
	}

	v, err := cache.Get(_const.CacheKeyConfig + _const.CacheKeyPayTimeout)
	if err != nil {
		// 默认30分钟
		return time.Duration(_const.PayTimeout) * time.Minute
	}
	var payTimeout int
	payTimeout, err = strconv.Atoi(v)
	if err != nil {
		// 默认30分钟
		return time.Duration(_const.PayTimeout) * time.Minute
	}
	// 默认30分钟
	return time.Duration(payTimeout) * time.Minute
}
