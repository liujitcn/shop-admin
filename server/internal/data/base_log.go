package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"time"
)

type BaseLogCondition struct {
	Id               int64      `search:"type:eq;column:id"`
	Operation        string     `search:"type:contains;column:operation"`
	StatusCode       int32      `search:"type:eq;column:status_code"`
	RequestStartTime *time.Time `search:"type:gte;column:request_time"`
	RequestEndTime   *time.Time `search:"type:lte;column:request_time"`
}

type BaseLogRepo interface {
	baseRepo.BaseRepo[models.BaseLog, BaseLogCondition]
}

type baseLogRepo struct {
	baseRepo.BaseRepo[models.BaseLog, BaseLogCondition]
	data *genData.Data
}

func NewBaseLogRepo(data *genData.Data) BaseLogRepo {
	base := baseRepo.NewBaseRepo[models.BaseLog, BaseLogCondition](
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
	)
	return &baseLogRepo{
		BaseRepo: base,
		data:     data,
	}
}
