package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ShopBannerCondition struct {
	Id     int64   `query:"type:eq;column:id"`
	Ids    []int64 `query:"type:in;column:id"`
	Site   int32   `query:"type:eq;column:site"`
	Type   int32   `query:"type:eq;column:type"`
	Status int32   `query:"type:eq;column:status"`
}

type ShopBannerRepo interface {
	genRepo.BaseRepo[models.ShopBanner, ShopBannerCondition]
}

type shopBannerRepo struct {
	genRepo.BaseRepo[models.ShopBanner, ShopBannerCondition]
	data *genData.Data
}

func NewShopBannerRepo(data *genData.Data) ShopBannerRepo {
	base := genRepo.NewBaseRepo[models.ShopBanner, ShopBannerCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopBanner.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).ShopBanner.ID
		},
		func(entity *models.ShopBanner) int64 {
			return entity.ID
		},
		new(models.ShopBanner),
		100,
	)
	return &shopBannerRepo{
		BaseRepo: base,
		data:     data,
	}
}
