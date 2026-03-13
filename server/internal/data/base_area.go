package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseAreaCondition struct {
	Id       int64   `search:"type:eq;column:id"`
	Ids      []int64 `search:"type:in;column:id"`
	ParentId int64   `search:"type:eq;column:parent_id"`
	Name     string  `search:"type:contains;column:name"`
	Code     string
}

type BaseAreaRepo interface {
	baseRepo.BaseRepo[models.BaseArea, BaseAreaCondition]
}

type baseAreaRepo struct {
	baseRepo.BaseRepo[models.BaseArea, BaseAreaCondition]
	data *genData.Data
}

func NewBaseAreaRepo(data *genData.Data) BaseAreaRepo {
	base := baseRepo.NewBaseRepo[models.BaseArea, BaseAreaCondition](
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
	)
	return &baseAreaRepo{
		BaseRepo: base,
		data:     data,
	}
}
