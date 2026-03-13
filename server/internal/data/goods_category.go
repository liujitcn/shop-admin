package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsCategoryCondition struct {
	Id            int64   `search:"type:eq;column:id"`
	Ids           []int64 `search:"type:in;column:id"`
	Name          string  `search:"type:contains;column:name"`
	ParentId      *int64  `search:"type:eq;column:parent_id"`
	Status        int32   `search:"type:eq;column:status"`
	ParentIDOrder string  `search:"type:order;column:parent_id"`
	SortOrder     string  `search:"type:order;column:sort"`
}

type GoodsCategoryRepo interface {
	baseRepo.BaseRepo[models.GoodsCategory, GoodsCategoryCondition]
}

type goodsCategoryRepo struct {
	baseRepo.BaseRepo[models.GoodsCategory, GoodsCategoryCondition]
	data *genData.Data
}

func NewGoodsCategoryRepo(data *genData.Data) GoodsCategoryRepo {
	base := baseRepo.NewBaseRepo[models.GoodsCategory, GoodsCategoryCondition](
		func(ctx context.Context) gen.Dao {
			dao := data.Query(ctx).GoodsCategory.WithContext(ctx).DO
			return &dao
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsCategory.ID
		},
		func(entity *models.GoodsCategory) int64 {
			return entity.ID
		},
		new(models.GoodsCategory),
	)
	return &goodsCategoryRepo{
		BaseRepo: base,
		data:     data,
	}
}
