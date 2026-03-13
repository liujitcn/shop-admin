package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsSpecCondition struct {
	Id      int64   `search:"type:eq;column:id"`
	Ids     []int64 `search:"type:in;column:id"`
	GoodsId int64   `search:"type:eq;column:goods_id"`
	Name    string  `search:"type:contains;column:name"`
	Status  int32   `search:"type:eq;column:status"`
}

type GoodsSpecRepo interface {
	baseRepo.BaseRepo[models.GoodsSpec, GoodsSpecCondition]
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
}

type goodsSpecRepo struct {
	baseRepo.BaseRepo[models.GoodsSpec, GoodsSpecCondition]
	data *genData.Data
}

func NewGoodsSpecRepo(data *genData.Data) GoodsSpecRepo {
	base := baseRepo.NewBaseRepo[models.GoodsSpec, GoodsSpecCondition](
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
