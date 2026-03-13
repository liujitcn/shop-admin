package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseRoleCondition struct {
	Id     int64   `search:"type:eq;column:id"`
	Ids    []int64 `search:"type:in;column:id"`
	Status int32   `search:"type:eq;column:status"`
	Name   string  `search:"type:contains;column:name"`
	Code   string  `search:"type:contains;column:code"`
}

type BaseRoleRepo interface {
	baseRepo.BaseRepo[models.BaseRole, BaseRoleCondition]
}

type baseRoleRepo struct {
	baseRepo.BaseRepo[models.BaseRole, BaseRoleCondition]
	data *genData.Data
}

func NewBaseRoleRepo(data *genData.Data) BaseRoleRepo {
	base := baseRepo.NewBaseRepo[models.BaseRole, BaseRoleCondition](
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
	)
	return &baseRoleRepo{
		BaseRepo: base,
		data:     data,
	}
}
