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
	Id     int64   `query:"type:eq;column:id"`
	Ids    []int64 `query:"type:in;column:id"`
	Status int32   `query:"type:eq;column:status"`
	Title  string  `query:"type:contains;column:title"` // 热门推荐标题
	Desc   string  `query:"type:contains;column:desc"`  // 热门推荐描述
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
