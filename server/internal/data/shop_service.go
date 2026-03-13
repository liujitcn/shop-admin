package data

import (
	"context"

	genData "github.com/liujitcn/shop-gorm-gen/data"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// ShopServiceCondition 商城服务查询条件
type ShopServiceCondition struct {
	Id     int64   `search:"type:eq;column:id"`
	Ids    []int64 `search:"type:in;column:id"`
	Label  string  `search:"type:contains;column:label"`
	Status int32   `search:"type:eq;column:status"`
}

// ShopServiceRepo 商城服务数据接口
type ShopServiceRepo interface {
	baseRepo.BaseRepo[models.ShopService, ShopServiceCondition]
}

// shopServiceRepo 商城服务数据实现
type shopServiceRepo struct {
	baseRepo.BaseRepo[models.ShopService, ShopServiceCondition]
	data *genData.Data
}

// NewShopServiceRepo 创建商城服务数据实例
func NewShopServiceRepo(data *genData.Data) ShopServiceRepo {
	base := baseRepo.NewBaseRepo[models.ShopService, ShopServiceCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopService.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).ShopService.ID
		},
		func(entity *models.ShopService) int64 {
			return entity.ID
		},
		new(models.ShopService),
	)
	return &shopServiceRepo{
		BaseRepo: base,
		data:     data,
	}
}
