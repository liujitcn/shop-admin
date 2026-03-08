package data

import (
	"context"

	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsCategoryCondition struct {
	Id            int64   `query:"type:eq;column:id"`
	Ids           []int64 `query:"type:in;column:id"`
	Name          string  `query:"type:contains;column:name"`
	ParentId      *int64  `query:"type:eq;column:parent_id"`
	Status        int32   `query:"type:eq;column:status"`
	ParentIDOrder string  `query:"type:order;column:parent_id"`
	SortOrder     string  `query:"type:order;column:sort"`
}

type GoodsCategoryRepo interface {
	genRepo.BaseRepo[models.GoodsCategory, GoodsCategoryCondition]
}

type goodsCategoryRepo struct {
	genRepo.BaseRepo[models.GoodsCategory, GoodsCategoryCondition]
	data *genData.Data
}

func NewGoodsCategoryRepo(data *genData.Data) GoodsCategoryRepo {
	base := genRepo.NewBaseRepo[models.GoodsCategory, GoodsCategoryCondition](
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
		100,
	)
	return &goodsCategoryRepo{
		BaseRepo: base,
		data:     data,
	}
}
