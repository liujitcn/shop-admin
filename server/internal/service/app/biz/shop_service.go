package biz

import (
	"context"

	"github.com/liujitcn/shop-admin/server/api/gen/go/app"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type ShopServiceCase struct {
	data.ShopServiceRepo
}

// NewShopServiceCase new a ShopService use case.
func NewShopServiceCase(shopServiceRepo data.ShopServiceRepo) *ShopServiceCase {
	return &ShopServiceCase{
		ShopServiceRepo: shopServiceRepo,
	}
}

func (c *ShopServiceCase) List(ctx context.Context, condition *data.ShopServiceCondition) ([]*models.ShopService, error) {
	return c.FindAll(ctx, condition)
}

func (c *ShopServiceCase) ConvertToProto(ctx context.Context, item *models.ShopService) *app.ShopService {
	return &app.ShopService{
		Label: item.Label,
		Value: item.Value,
	}
}
