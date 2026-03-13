//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	baseConfigs "github.com/liujitcn/shop-base/server/configs"
	baseCore "github.com/liujitcn/shop-base/server/core"
	baseData "github.com/liujitcn/shop-base/server/data"
	baseService "github.com/liujitcn/shop-base/server/service"

	"github.com/liujitcn/kratos-kit/bootstrap"
	"github.com/liujitcn/shop-admin/server/internal/configs"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-admin/server/internal/middleware"
	"github.com/liujitcn/shop-admin/server/internal/server"
	"github.com/liujitcn/shop-admin/server/internal/service"
)

// initApp init kratos application.
func initApp(*bootstrap.Context) (*kratos.App, func(), error) {
	panic(wire.Build(
		baseConfigs.ProviderSet,
		baseCore.ProviderSet,
		baseData.ProviderSet,
		baseService.ProviderSet,
		configs.ProviderSet,
		data.ProviderSet,
		middleware.ProviderSet,
		server.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}
