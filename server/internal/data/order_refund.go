package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderRefundCondition struct {
	OrderId     int64 `search:"type:eq;column:order_id"`
	Status      int32 `search:"type:eq;column:status"`
	SuccessTime string
}

type OrderRefundRepo interface {
	baseRepo.BaseRepo[models.OrderRefund, OrderRefundCondition]
}

type orderRefundRepo struct {
	baseRepo.BaseRepo[models.OrderRefund, OrderRefundCondition]
	data *genData.Data
}

func NewOrderRefundRepo(data *genData.Data) OrderRefundRepo {
	base := baseRepo.NewBaseRepo[models.OrderRefund, OrderRefundCondition](
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
	)
	return &orderRefundRepo{
		BaseRepo: base,
		data:     data,
	}
}
