package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsSpecCondition struct {
	Id      int64
	Ids     []int64
	GoodsId int64
	Name    string
	Status  int32
}

type GoodsSpecRepo interface {
	genRepo.BaseRepo[models.GoodsSpec, GoodsSpecCondition]
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
}

type goodsSpecRepo struct {
	genRepo.BaseRepo[models.GoodsSpec, GoodsSpecCondition]
	data *genData.Data
}

func NewGoodsSpecRepo(data *genData.Data) GoodsSpecRepo {
	base := genRepo.NewBaseRepo[models.GoodsSpec, GoodsSpecCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).GoodsSpec.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsSpec.ID
		},
		func(entity *models.GoodsSpec) int64 {
			return entity.ID
		},
		new(models.GoodsSpec),
		100,
	)
	return &goodsSpecRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *goodsSpecRepo) DeleteByGoodsId(ctx context.Context, goodsId int64) error {
	q := r.data.Query(ctx).GoodsSpec
	_, err := q.WithContext(ctx).Where(q.GoodsID.Eq(goodsId)).Delete()
	return err
}
