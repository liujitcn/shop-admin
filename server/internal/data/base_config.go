package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseConfigCondition struct {
	Id     int64   `search:"type:eq;column:id"`
	Ids    []int64 `search:"type:in;column:id"`
	Site   int32   `search:"type:eq;column:site"`
	Name   string  `search:"type:contains;column:name"`
	Type   int32   `search:"type:eq;column:type"`
	Key    string  `search:"type:contains;column:key"`
	Status int32   `search:"type:eq;column:status"`
}

type BaseConfigRepo interface {
	baseRepo.BaseRepo[models.BaseConfig, BaseConfigCondition]
}

type baseConfigRepo struct {
	baseRepo.BaseRepo[models.BaseConfig, BaseConfigCondition]
	data *genData.Data
}

func NewBaseConfigRepo(data *genData.Data) BaseConfigRepo {
	base := baseRepo.NewBaseRepo[models.BaseConfig, BaseConfigCondition](
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
	)
	return &baseConfigRepo{
		BaseRepo: base,
		data:     data,
	}
}
