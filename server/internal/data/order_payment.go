package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderPaymentCondition struct {
	OrderId     int64 `query:"type:eq;column:order_id"`
	Status      int32 `query:"type:eq;column:status"`
	SuccessTime string
}

type OrderPaymentRepo interface {
	genRepo.BaseRepo[models.OrderPayment, OrderPaymentCondition]
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderPaymentRepo struct {
	genRepo.BaseRepo[models.OrderPayment, OrderPaymentCondition]
	data *genData.Data
}

func NewOrderPaymentRepo(data *genData.Data) OrderPaymentRepo {
	base := genRepo.NewBaseRepo[models.OrderPayment, OrderPaymentCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderPayment.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderPayment.ID
		},
		func(entity *models.OrderPayment) int64 {
			return entity.ID
		},
		new(models.OrderPayment),
		100,
	)
	return &orderPaymentRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *orderPaymentRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderPayment
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
