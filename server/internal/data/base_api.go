package data

import (
	"context"

	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseApiCondition struct {
	Id          int64   `query:"type:eq;column:id"`
	Ids         []int64 `query:"type:in;column:id"`
	ServiceName string  `query:"type:contains;column:service_name"` // 服务名
	ServiceDesc string  `query:"type:contains;column:service_desc"` // 服务描述
	Desc        string  `query:"type:contains;column:desc"`         // 描述
	Path        string  `query:"type:contains;column:path"`         // 请求地址
}

type BaseApiRepo interface {
	genRepo.BaseRepo[models.BaseAPI, BaseApiCondition]
}

type baseApiRepo struct {
	genRepo.BaseRepo[models.BaseAPI, BaseApiCondition]
	data *genData.Data
}

func NewBaseApiRepo(data *genData.Data) BaseApiRepo {
	base := genRepo.NewBaseRepo[models.BaseAPI, BaseApiCondition](
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
		100,
	)
	return &baseApiRepo{
		BaseRepo: base,
		data:     data,
	}
}
