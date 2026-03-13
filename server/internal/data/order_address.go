package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderAddressCondition struct {
	OrderId int64 `search:"type:eq;column:order_id"`
}

type OrderAddressRepo interface {
	baseRepo.BaseRepo[models.OrderAddress, OrderAddressCondition]
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderAddressRepo struct {
	baseRepo.BaseRepo[models.OrderAddress, OrderAddressCondition]
	data *genData.Data
}

func NewOrderAddressRepo(data *genData.Data) OrderAddressRepo {
	base := baseRepo.NewBaseRepo[models.OrderAddress, OrderAddressCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderAddress.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderAddress.ID
		},
		func(entity *models.OrderAddress) int64 {
			return entity.ID
		},
		new(models.OrderAddress),
	)
	return &orderAddressRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *orderAddressRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderAddress
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
