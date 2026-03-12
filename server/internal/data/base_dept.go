package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseDeptCondition struct {
	Id       int64  `query:"type:eq;column:id"`
	ParentId *int64 `query:"type:eq;column:parent_id"`
	Status   int32  `query:"type:eq;column:status"`
}

type BaseDeptRepo interface {
	genRepo.BaseRepo[models.BaseDept, BaseDeptCondition]
}

type baseDeptRepo struct {
	genRepo.BaseRepo[models.BaseDept, BaseDeptCondition]
	data *genData.Data
}

func NewBaseDeptRepo(data *genData.Data) BaseDeptRepo {
	base := genRepo.NewBaseRepo[models.BaseDept, BaseDeptCondition](
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
		100,
	)
	return &baseDeptRepo{
		BaseRepo: base,
		data:     data,
	}
}
