package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"time"
)

type BaseLogCondition struct {
	Id               int64
	Operation        string
	StatusCode       int32
	RequestStartTime *time.Time
	RequestEndTime   *time.Time
}

type BaseLogRepo interface {
	genRepo.BaseRepo[models.BaseLog, BaseLogCondition]
}

type baseLogRepo struct {
	genRepo.BaseRepo[models.BaseLog, BaseLogCondition]
	data *genData.Data
}

func NewBaseLogRepo(data *genData.Data) BaseLogRepo {
	base := genRepo.NewBaseRepo[models.BaseLog, BaseLogCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseLog.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseLog.ID
		},
		func(entity *models.BaseLog) int64 {
			return entity.ID
		},
		new(models.BaseLog),
		100,
	)
	return &baseLogRepo{
		BaseRepo: base,
		data:     data,
	}
}
