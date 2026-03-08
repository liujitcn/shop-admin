package data

import (
	"context"

	genData "github.com/liujitcn/shop-gorm-gen/data"

	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

// ShopServiceCondition 商城服务查询条件
type ShopServiceCondition struct {
	Id     int64
	Ids    []int64
	Label  string
	Status int32
}

// ShopServiceRepo 商城服务数据接口
type ShopServiceRepo interface {
	genRepo.BaseRepo[models.ShopService, ShopServiceCondition]
}

// shopServiceRepo 商城服务数据实现
type shopServiceRepo struct {
	genRepo.BaseRepo[models.ShopService, ShopServiceCondition]
	data *genData.Data
}

// NewShopServiceRepo 创建商城服务数据实例
func NewShopServiceRepo(data *genData.Data) ShopServiceRepo {
	base := genRepo.NewBaseRepo[models.ShopService, ShopServiceCondition](
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
		100,
	)
	return &shopServiceRepo{
		BaseRepo: base,
		data:     data,
	}
}
