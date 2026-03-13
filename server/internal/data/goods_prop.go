package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsPropCondition struct {
	Id      int64   `search:"type:eq;column:id"`
	Ids     []int64 `search:"type:in;column:id"`
	GoodsId int64   `search:"type:eq;column:goods_id"`
	Label   string  `search:"type:contains;column:label"`
	Status  int32   `search:"type:eq;column:status"`
}

type GoodsPropRepo interface {
	baseRepo.BaseRepo[models.GoodsProp, GoodsPropCondition]
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
}

type goodsPropRepo struct {
	baseRepo.BaseRepo[models.GoodsProp, GoodsPropCondition]
	data *genData.Data
}

func NewGoodsPropRepo(data *genData.Data) GoodsPropRepo {
	base := baseRepo.NewBaseRepo[models.GoodsProp, GoodsPropCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).GoodsProp.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsProp.ID
		},
		func(entity *models.GoodsProp) int64 {
			return entity.ID
		},
		new(models.GoodsProp),
	)
	return &goodsPropRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *goodsPropRepo) DeleteByGoodsId(ctx context.Context, goodsId int64) error {
	q := r.data.Query(ctx).GoodsProp
	_, err := q.WithContext(ctx).Where(q.GoodsID.Eq(goodsId)).Delete()
	return err
}
