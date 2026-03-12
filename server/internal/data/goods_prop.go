package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsPropCondition struct {
	Id      int64   `query:"type:eq;column:id"`
	Ids     []int64 `query:"type:in;column:id"`
	GoodsId int64   `query:"type:eq;column:goods_id"`
	Label   string  `query:"type:contains;column:label"`
	Status  int32   `query:"type:eq;column:status"`
}

type GoodsPropRepo interface {
	genRepo.BaseRepo[models.GoodsProp, GoodsPropCondition]
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
}

type goodsPropRepo struct {
	genRepo.BaseRepo[models.GoodsProp, GoodsPropCondition]
	data *genData.Data
}

func NewGoodsPropRepo(data *genData.Data) GoodsPropRepo {
	base := genRepo.NewBaseRepo[models.GoodsProp, GoodsPropCondition](
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
		100,
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
