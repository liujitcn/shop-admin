package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseMenuCondition struct {
	Id       int64   `search:"type:eq;column:id"`
	ParentId *int64  `search:"type:eq;column:parent_id"`
	Ids      []int64 `search:"type:in;column:id"`
	Status   int32   `search:"type:eq;column:status"`
	Types    []int32 `search:"type:in;column:type"`
}

type BaseMenuRepo interface {
	baseRepo.BaseRepo[models.BaseMenu, BaseMenuCondition]
}

type baseMenuRepo struct {
	baseRepo.BaseRepo[models.BaseMenu, BaseMenuCondition]
	data *genData.Data
}

func NewBaseMenuRepo(data *genData.Data) BaseMenuRepo {
	base := baseRepo.NewBaseRepo[models.BaseMenu, BaseMenuCondition](
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
	)
	return &baseMenuRepo{
		BaseRepo: base,
		data:     data,
	}
}
