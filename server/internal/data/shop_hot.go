package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ShopHotCondition struct {
	Id     int64
	Ids    []int64
	Status int32
	Title  string // 热门推荐标题
	Desc   string // 热门推荐描述
}

type ShopHotRepo interface {
	genRepo.BaseRepo[models.ShopHot, ShopHotCondition]
}

type shopHotRepo struct {
	genRepo.BaseRepo[models.ShopHot, ShopHotCondition]
	data *genData.Data
}

func NewShopHotRepo(data *genData.Data) ShopHotRepo {
	base := genRepo.NewBaseRepo[models.ShopHot, ShopHotCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).ShopHot.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).ShopHot.ID
		},
		func(entity *models.ShopHot) int64 {
			return entity.ID
		},
		new(models.ShopHot),
		100,
	)
	return &shopHotRepo{
		BaseRepo: base,
		data:     data,
	}
}
