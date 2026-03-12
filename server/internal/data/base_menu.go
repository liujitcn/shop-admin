package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseMenuCondition struct {
	Id       int64   `query:"type:eq;column:id"`
	ParentId *int64  `query:"type:eq;column:parent_id"`
	Ids      []int64 `query:"type:in;column:id"`
	Status   int32   `query:"type:eq;column:status"`
	Types    []int32 `query:"type:in;column:type"`
}

type BaseMenuRepo interface {
	genRepo.BaseRepo[models.BaseMenu, BaseMenuCondition]
}

type baseMenuRepo struct {
	genRepo.BaseRepo[models.BaseMenu, BaseMenuCondition]
	data *genData.Data
}

func NewBaseMenuRepo(data *genData.Data) BaseMenuRepo {
	base := genRepo.NewBaseRepo[models.BaseMenu, BaseMenuCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseMenu.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseMenu.ID
		},
		func(entity *models.BaseMenu) int64 {
			return entity.ID
		},
		new(models.BaseMenu),
		100,
	)
	return &baseMenuRepo{
		BaseRepo: base,
		data:     data,
	}
}
