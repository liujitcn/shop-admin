package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ShopHotItemCondition struct {
	Id     int64  `search:"type:eq;column:id"`
	HotId  int64  `search:"type:eq;column:hot_id"`
	Status int32  `search:"type:eq;column:status"`
	Title  string `search:"type:contains;column:title"`
}

type ShopHotItemRepo interface {
	baseRepo.BaseRepo[models.ShopHotItem, ShopHotItemCondition]
}

type shopHotItemRepo struct {
	baseRepo.BaseRepo[models.ShopHotItem, ShopHotItemCondition]
	data *genData.Data
}

func NewShopHotItemRepo(data *genData.Data) ShopHotItemRepo {
	base := baseRepo.NewBaseRepo[models.ShopHotItem, ShopHotItemCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopHotItem.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).ShopHotItem.ID
		},
		func(entity *models.ShopHotItem) int64 {
			return entity.ID
		},
		new(models.ShopHotItem),
	)
	return &shopHotItemRepo{
		BaseRepo: base,
		data:     data,
	}
}
