package data

import (
	"context"

	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseRoleCondition struct {
	Id     int64   `query:"type:eq;column:id"`
	Ids    []int64 `query:"type:in;column:id"`
	Status int32   `query:"type:eq;column:status"`
	Name   string  `query:"type:contains;column:name"`
	Code   string  `query:"type:contains;column:code"`
}

type BaseRoleRepo interface {
	genRepo.BaseRepo[models.BaseRole, BaseRoleCondition]
}

type baseRoleRepo struct {
	genRepo.BaseRepo[models.BaseRole, BaseRoleCondition]
	data *genData.Data
}

func NewBaseRoleRepo(data *genData.Data) BaseRoleRepo {
	base := genRepo.NewBaseRepo[models.BaseRole, BaseRoleCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseRole.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseRole.ID
		},
		func(entity *models.BaseRole) int64 {
			return entity.ID
		},
		new(models.BaseRole),
		100,
	)
	return &baseRoleRepo{
		BaseRepo: base,
		data:     data,
	}
}
