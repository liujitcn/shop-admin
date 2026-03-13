package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ShopBannerCondition struct {
	Id     int64   `search:"type:eq;column:id"`
	Ids    []int64 `search:"type:in;column:id"`
	Site   int32   `search:"type:eq;column:site"`
	Type   int32   `search:"type:eq;column:type"`
	Status int32   `search:"type:eq;column:status"`
}

type ShopBannerRepo interface {
	baseRepo.BaseRepo[models.ShopBanner, ShopBannerCondition]
}

type shopBannerRepo struct {
	baseRepo.BaseRepo[models.ShopBanner, ShopBannerCondition]
	data *genData.Data
}

func NewShopBannerRepo(data *genData.Data) ShopBannerRepo {
	base := baseRepo.NewBaseRepo[models.ShopBanner, ShopBannerCondition](
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
	)
	return &shopBannerRepo{
		BaseRepo: base,
		data:     data,
	}
}
