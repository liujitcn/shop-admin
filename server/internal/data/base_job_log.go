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

type BaseJobLogCondition struct {
	Id               int64      `search:"type:eq;column:id"`
	JobId            int64      `search:"type:eq;column:job_id"`
	Status           int32      `search:"type:eq;column:status"`
	ExecuteStartTime *time.Time `search:"type:gte;column:execute_time"`
	ExecuteEndTime   *time.Time `search:"type:lte;column:execute_time"`
}

type BaseJobLogRepo interface {
	baseRepo.BaseRepo[models.BaseJobLog, BaseJobLogCondition]
}

type baseJobLogRepo struct {
	baseRepo.BaseRepo[models.BaseJobLog, BaseJobLogCondition]
	data *genData.Data
}

func NewBaseJobLogRepo(data *genData.Data) BaseJobLogRepo {
	base := baseRepo.NewBaseRepo[models.BaseJobLog, BaseJobLogCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseJobLog.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseJobLog.ID
		},
		func(entity *models.BaseJobLog) int64 {
			return entity.ID
		},
		new(models.BaseJobLog),
	)
	return &baseJobLogRepo{
		BaseRepo: base,
		data:     data,
	}
}
