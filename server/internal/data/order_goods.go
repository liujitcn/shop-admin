package data

import (
	"context"
	"time"

	"github.com/liujitcn/shop-admin/server/internal/data/dto"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderGoodsCondition struct {
	OrderId  int64
	OrderIds []int64
}

type OrderGoodsRepo interface {
	genRepo.BaseRepo[models.OrderGoods, OrderGoodsCondition]
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
	OrderGoodsStatusSummary(ctx context.Context, startCreatedAt, endCreatedAt *time.Time) ([]*dto.OrderGoodsStatusSummary, error)
	OrderGoodsSummary(ctx context.Context, top int64, startCreatedAt, endCreatedAt *time.Time) ([]*dto.OrderGoodsSummary, error)
}

type orderGoodsRepo struct {
	genRepo.BaseRepo[models.OrderGoods, OrderGoodsCondition]
	data *genData.Data
}

func NewOrderGoodsRepo(data *genData.Data) OrderGoodsRepo {
	base := genRepo.NewBaseRepo[models.OrderGoods, OrderGoodsCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderGoods.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderGoods.ID
		},
		func(entity *models.OrderGoods) int64 {
			return entity.ID
		},
		new(models.OrderGoods),
		100,
	)
	return &orderGoodsRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *orderGoodsRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderGoods
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}

func (r *orderGoodsRepo) OrderGoodsStatusSummary(ctx context.Context, startCreatedAt, endCreatedAt *time.Time) ([]*dto.OrderGoodsStatusSummary, error) {
	order := r.data.Query(ctx).Order
	goods := r.data.Query(ctx).Goods
	goodsCategory := r.data.Query(ctx).GoodsCategory
	m := r.data.Query(ctx).OrderGoods
	q := m.WithContext(ctx)
	q = q.Join(order, order.ID.EqCol(m.OrderID))
	q = q.Join(goods, goods.ID.EqCol(m.GoodsID))
	q = q.Join(goodsCategory, goodsCategory.ID.EqCol(goods.CategoryID))

	if startCreatedAt != nil {
		q = q.Where(order.CreatedAt.Gte(*startCreatedAt))
	}
	if endCreatedAt != nil {
		q = q.Where(order.CreatedAt.Lt(*endCreatedAt))
	}
	results := make([]*dto.OrderGoodsStatusSummary, 0)

	q.Select(goodsCategory.ParentID.As("category_id"), order.Status.As("status"), m.Num.Sum().As("goods_count")).Group(goodsCategory.ParentID, order.Status)
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *orderGoodsRepo) OrderGoodsSummary(ctx context.Context, top int64, startCreatedAt, endCreatedAt *time.Time) ([]*dto.OrderGoodsSummary, error) {
	order := r.data.Query(ctx).Order
	m := r.data.Query(ctx).OrderGoods
	q := m.WithContext(ctx)
	q = q.Join(order, order.ID.EqCol(m.OrderID))

	if startCreatedAt != nil {
		q = q.Where(order.CreatedAt.Gte(*startCreatedAt))
	}
	if endCreatedAt != nil {
		q = q.Where(order.CreatedAt.Lt(*endCreatedAt))
	}
	results := make([]*dto.OrderGoodsSummary, 0)

	q.Select(m.GoodsID.As("goods_id"), m.Num.Sum().As("goods_count")).Group(m.GoodsID).Order(m.Num.Sum().Desc()).Limit(int(top))
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
