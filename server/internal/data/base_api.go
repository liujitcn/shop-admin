package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseApiCondition struct {
	Id          int64   `search:"type:eq;column:id"`
	Ids         []int64 `search:"type:in;column:id"`
	ServiceName string  `search:"type:contains;column:service_name"` // 服务名
	ServiceDesc string  `search:"type:contains;column:service_desc"` // 服务描述
	Desc        string  `search:"type:contains;column:desc"`         // 描述
	Path        string  `search:"type:contains;column:path"`         // 请求地址
}

type BaseApiRepo interface {
	baseRepo.BaseRepo[models.BaseAPI, BaseApiCondition]
}

type baseApiRepo struct {
	baseRepo.BaseRepo[models.BaseAPI, BaseApiCondition]
	data *genData.Data
}

func NewBaseApiRepo(data *genData.Data) BaseApiRepo {
	base := baseRepo.NewBaseRepo[models.BaseAPI, BaseApiCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseAPI.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseAPI.ID
		},
		func(entity *models.BaseAPI) int64 {
			return entity.ID
		},
		new(models.BaseAPI),
	)
	return &baseApiRepo{
		BaseRepo: base,
		data:     data,
	}
}
