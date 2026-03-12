package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ShopHotItemCondition struct {
	Id     int64  `query:"type:eq;column:id"`
	HotId  int64  `query:"type:eq;column:hot_id"`
	Status int32  `query:"type:eq;column:status"`
	Title  string `query:"type:contains;column:title"`
}

type ShopHotItemRepo interface {
	genRepo.BaseRepo[models.ShopHotItem, ShopHotItemCondition]
}

type shopHotItemRepo struct {
	genRepo.BaseRepo[models.ShopHotItem, ShopHotItemCondition]
	data *genData.Data
}

func NewShopHotItemRepo(data *genData.Data) ShopHotItemRepo {
	base := genRepo.NewBaseRepo[models.ShopHotItem, ShopHotItemCondition](
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
		100,
	)
	return &shopHotItemRepo{
		BaseRepo: base,
		data:     data,
	}
}
