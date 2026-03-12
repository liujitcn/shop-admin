package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderRefundCondition struct {
	OrderId     int64 `query:"type:eq;column:order_id"`
	Status      int32 `query:"type:eq;column:status"`
	SuccessTime string
}

type OrderRefundRepo interface {
	genRepo.BaseRepo[models.OrderRefund, OrderRefundCondition]
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderRefundRepo struct {
	genRepo.BaseRepo[models.OrderRefund, OrderRefundCondition]
	data *genData.Data
}

func NewOrderRefundRepo(data *genData.Data) OrderRefundRepo {
	base := genRepo.NewBaseRepo[models.OrderRefund, OrderRefundCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderRefund.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderRefund.ID
		},
		func(entity *models.OrderRefund) int64 {
			return entity.ID
		},
		new(models.OrderRefund),
		100,
	)
	return &orderRefundRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *orderRefundRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderRefund
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
