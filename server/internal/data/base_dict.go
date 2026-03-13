package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseDictCondition struct {
	Id     int64    `search:"type:eq;column:id"`
	Status int32    `search:"type:eq;column:status"`
	Name   string   `search:"type:contains;column:name"` // 字典名称
	Code   string   `search:"type:contains;column:code"` // 字典代码
	Codes  []string `search:"type:in;column:code"`       // 字典代码
}

type BaseDictRepo interface {
	baseRepo.BaseRepo[models.BaseDict, BaseDictCondition]
}

type baseDictRepo struct {
	baseRepo.BaseRepo[models.BaseDict, BaseDictCondition]
	data *genData.Data
}

func NewBaseDictRepo(data *genData.Data) BaseDictRepo {
	base := baseRepo.NewBaseRepo[models.BaseDict, BaseDictCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseDict.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseDict.ID
		},
		func(entity *models.BaseDict) int64 {
			return entity.ID
		},
		new(models.BaseDict),
	)
	return &baseDictRepo{
		BaseRepo: base,
		data:     data,
	}
}
