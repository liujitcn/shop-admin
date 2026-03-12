package data

import (
	"context"

	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseAreaCondition struct {
	Id       int64   `query:"type:eq;column:id"`
	Ids      []int64 `query:"type:in;column:id"`
	ParentId int64   `query:"type:eq;column:parent_id"`
	Name     string  `query:"type:contains;column:name"`
	Code     string
}

type BaseAreaRepo interface {
	genRepo.BaseRepo[models.BaseArea, BaseAreaCondition]
}

type baseAreaRepo struct {
	genRepo.BaseRepo[models.BaseArea, BaseAreaCondition]
	data *genData.Data
}

func NewBaseAreaRepo(data *genData.Data) BaseAreaRepo {
	base := genRepo.NewBaseRepo[models.BaseArea, BaseAreaCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseArea.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseArea.ID
		},
		func(entity *models.BaseArea) int64 {
			return entity.ID
		},
		new(models.BaseArea),
		100,
	)
	return &baseAreaRepo{
		BaseRepo: base,
		data:     data,
	}
}
