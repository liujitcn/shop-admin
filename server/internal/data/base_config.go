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
	Id     int64
	Ids    []int64
	Site   int32
	Name   string
	Type   int32
	Key    string
	Status int32
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
