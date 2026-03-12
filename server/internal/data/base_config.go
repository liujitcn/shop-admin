package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseConfigCondition struct {
	Id     int64   `query:"type:eq;column:id"`
	Ids    []int64 `query:"type:in;column:id"`
	Site   int32   `query:"type:eq;column:site"`
	Name   string  `query:"type:contains;column:name"`
	Type   int32   `query:"type:eq;column:type"`
	Key    string  `query:"type:contains;column:key"`
	Status int32   `query:"type:eq;column:status"`
}

type BaseConfigRepo interface {
	genRepo.BaseRepo[models.BaseConfig, BaseConfigCondition]
}

type baseConfigRepo struct {
	genRepo.BaseRepo[models.BaseConfig, BaseConfigCondition]
	data *genData.Data
}

func NewBaseConfigRepo(data *genData.Data) BaseConfigRepo {
	base := genRepo.NewBaseRepo[models.BaseConfig, BaseConfigCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseConfig.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseConfig.ID
		},
		func(entity *models.BaseConfig) int64 {
			return entity.ID
		},
		new(models.BaseConfig),
		100,
	)
	return &baseConfigRepo{
		BaseRepo: base,
		data:     data,
	}
}
