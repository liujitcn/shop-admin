package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderCancelCondition struct {
	OrderId int64 `search:"type:eq;column:order_id"`
}

type OrderCancelRepo interface {
	baseRepo.BaseRepo[models.OrderCancel, OrderCancelCondition]
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderCancelRepo struct {
	baseRepo.BaseRepo[models.OrderCancel, OrderCancelCondition]
	data *genData.Data
}

func NewOrderCancelRepo(data *genData.Data) OrderCancelRepo {
	base := baseRepo.NewBaseRepo[models.OrderCancel, OrderCancelCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderCancel.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderCancel.ID
		},
		func(entity *models.OrderCancel) int64 {
			return entity.ID
		},
		new(models.OrderCancel),
	)
	return &orderCancelRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *orderCancelRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderCancel
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
