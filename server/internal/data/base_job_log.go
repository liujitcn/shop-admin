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

type BaseJobLogCondition struct {
	Id               int64
	JobId            int64
	Status           int32
	ExecuteStartTime *time.Time
	ExecuteEndTime   *time.Time
}

type BaseJobLogRepo interface {
	genRepo.BaseRepo[models.BaseJobLog, BaseJobLogCondition]
}

type baseJobLogRepo struct {
	genRepo.BaseRepo[models.BaseJobLog, BaseJobLogCondition]
	data *genData.Data
}

func NewBaseJobLogRepo(data *genData.Data) BaseJobLogRepo {
	base := genRepo.NewBaseRepo[models.BaseJobLog, BaseJobLogCondition](
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
		100,
	)
	return &baseJobLogRepo{
		BaseRepo: base,
		data:     data,
	}
}
