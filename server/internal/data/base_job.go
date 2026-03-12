package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseJobCondition struct {
	Id           int64  `query:"type:eq;column:id"`
	Status       int32  `query:"type:eq;column:status"`
	Name         string `query:"type:contains;column:name"`
	InvokeTarget string `query:"type:contains;column:invoke_target"`
}

type BaseJobRepo interface {
	genRepo.BaseRepo[models.BaseJob, BaseJobCondition]
	CleanEntryId(ctx context.Context, entryId int32) error
}

type baseJobRepo struct {
	genRepo.BaseRepo[models.BaseJob, BaseJobCondition]
	data *genData.Data
}

func NewBaseJobRepo(data *genData.Data) BaseJobRepo {
	base := genRepo.NewBaseRepo[models.BaseJob, BaseJobCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseJob.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseJob.ID
		},
		func(entity *models.BaseJob) int64 {
			return entity.ID
		},
		new(models.BaseJob),
		100,
	)
	return &baseJobRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *baseJobRepo) CleanEntryId(ctx context.Context, entryId int32) error {
	var err error
	m := r.data.Query(ctx).BaseJob
	q := m.WithContext(ctx)
	if entryId > 0 {
		_, err = q.Where(m.EntryID.Eq(entryId)).UpdateColumn(m.EntryID, 0)
	} else {
		_, err = q.Where(m.EntryID.Gt(entryId)).Updates(&models.BaseJob{
			EntryID: 0,
		})
	}
	return err
}
