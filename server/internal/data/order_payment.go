package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderPaymentCondition struct {
	OrderId     int64 `search:"type:eq;column:order_id"`
	Status      int32 `search:"type:eq;column:status"`
	SuccessTime string
}

type OrderPaymentRepo interface {
	baseRepo.BaseRepo[models.OrderPayment, OrderPaymentCondition]
}

type orderPaymentRepo struct {
	baseRepo.BaseRepo[models.OrderPayment, OrderPaymentCondition]
	data *genData.Data
}

func NewOrderPaymentRepo(data *genData.Data) OrderPaymentRepo {
	base := baseRepo.NewBaseRepo[models.OrderPayment, OrderPaymentCondition](
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
	)
	return &orderPaymentRepo{
		BaseRepo: base,
		data:     data,
	}
}
