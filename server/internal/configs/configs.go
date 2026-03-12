package configs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	bootstrapConf "github.com/liujitcn/kratos-kit/api/gen/go/conf"
	"github.com/liujitcn/kratos-kit/bootstrap"
	"github.com/liujitcn/kratos-kit/sdk"
	"github.com/liujitcn/shop-admin/server/api/gen/go/conf"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
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

func ParseOss(ctx *bootstrap.Context) (*bootstrapConf.OSS, error) {
	cfg := ctx.GetConfig()
	if cfg == nil || cfg.GetOss() == nil {
		return nil, errors.New("config oss is nil")
	}
	return cfg.GetOss(), nil
}

func ParseData(ctx *bootstrap.Context) (*bootstrapConf.Data, error) {
	cfg := ctx.GetConfig()
	if cfg == nil || cfg.GetData() == nil {
		return nil, errors.New("config data is nil")
	}
	return cfg.GetData(), nil
}

func ParseDatabase(cfg *bootstrapConf.Data) *bootstrapConf.Data_Database {
	return cfg.GetDatabase()
}

func ParseQueue(cfg *bootstrapConf.Data) *bootstrapConf.Data_Queue {
	return cfg.GetQueue()
}

func ParseRedis(cfg *bootstrapConf.Data) *bootstrapConf.Data_Redis {
	return cfg.GetRedis()
}

func ParsePprof(ctx *bootstrap.Context) (*bootstrapConf.Pprof, error) {
	cfg := ctx.GetConfig()
	if cfg == nil || cfg.GetPprof() == nil {
		return nil, errors.New("config pprof is nil")
	}
	return cfg.GetPprof(), nil
}

func ParseAuthnJwt(ctx *bootstrap.Context) *bootstrapConf.Authentication_Jwt {
	cfg := ctx.GetConfig()
	if cfg == nil || cfg.GetAuthn() == nil || cfg.GetAuthn().GetJwt() == nil {
		return &bootstrapConf.Authentication_Jwt{
			Method: "HS256",
			Secret: "shop-admin",
		}
	}
	return cfg.GetAuthn().GetJwt()
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
