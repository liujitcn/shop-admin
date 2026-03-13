package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseDeptCondition struct {
	Id       int64  `search:"type:eq;column:id"`
	ParentId *int64 `search:"type:eq;column:parent_id"`
	Status   int32  `search:"type:eq;column:status"`
}

type BaseDeptRepo interface {
	baseRepo.BaseRepo[models.BaseDept, BaseDeptCondition]
}

type baseDeptRepo struct {
	baseRepo.BaseRepo[models.BaseDept, BaseDeptCondition]
	data *genData.Data
}

func NewBaseDeptRepo(data *genData.Data) BaseDeptRepo {
	base := baseRepo.NewBaseRepo[models.BaseDept, BaseDeptCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseDept.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseDept.ID
		},
		func(entity *models.BaseDept) int64 {
			return entity.ID
		},
		new(models.BaseDept),
	)
	return &baseDeptRepo{
		BaseRepo: base,
		data:     data,
	}
}
