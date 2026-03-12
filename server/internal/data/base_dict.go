package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseDictCondition struct {
	Id     int64    `query:"type:eq;column:id"`
	Status int32    `query:"type:eq;column:status"`
	Name   string   `query:"type:contains;column:name"` // 字典名称
	Code   string   `query:"type:contains;column:code"` // 字典代码
	Codes  []string `query:"type:in;column:code"`       // 字典代码
}

type BaseDictRepo interface {
	genRepo.BaseRepo[models.BaseDict, BaseDictCondition]
}

type baseDictRepo struct {
	genRepo.BaseRepo[models.BaseDict, BaseDictCondition]
	data *genData.Data
}

func NewBaseDictRepo(data *genData.Data) BaseDictRepo {
	base := genRepo.NewBaseRepo[models.BaseDict, BaseDictCondition](
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
		100,
	)
	return &baseDictRepo{
		BaseRepo: base,
		data:     data,
	}
}
