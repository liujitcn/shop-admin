//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-admin/server/internal/sdk"

	"github.com/go-kratos/kratos/v2"
	"github.com/liujitcn/shop-admin/server/internal/configs"
	"github.com/liujitcn/shop-admin/server/internal/server"
	"github.com/liujitcn/shop-admin/server/internal/service"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
)

// initApp init kratos application.
func initApp(*bootstrap.Context) (*kratos.App, func(), error) {
	panic(wire.Build(configs.ProviderSet, data.ProviderSet, sdk.ProviderSet, server.ProviderSet, service.ProviderSet, newApp))
}
