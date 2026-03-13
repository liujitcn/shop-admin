package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ShopHotCondition struct {
	Id     int64   `search:"type:eq;column:id"`
	Ids    []int64 `search:"type:in;column:id"`
	Status int32   `search:"type:eq;column:status"`
	Title  string  `search:"type:contains;column:title"` // 热门推荐标题
	Desc   string  `search:"type:contains;column:desc"`  // 热门推荐描述
}

type ShopHotRepo interface {
	baseRepo.BaseRepo[models.ShopHot, ShopHotCondition]
}

type shopHotRepo struct {
	baseRepo.BaseRepo[models.ShopHot, ShopHotCondition]
	data *genData.Data
}

func NewShopHotRepo(data *genData.Data) ShopHotRepo {
	base := baseRepo.NewBaseRepo[models.ShopHot, ShopHotCondition](
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
	)
	return &shopHotRepo{
		BaseRepo: base,
		data:     data,
	}
}
