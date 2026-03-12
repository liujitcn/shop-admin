package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderLogisticsCondition struct {
	OrderId int64 `query:"type:eq;column:order_id"`
}

type OrderLogisticsRepo interface {
	genRepo.BaseRepo[models.OrderLogistics, OrderLogisticsCondition]
	DeleteByOrderIds(ctx context.Context, orderIds []int64) error
}

type orderLogisticsRepo struct {
	genRepo.BaseRepo[models.OrderLogistics, OrderLogisticsCondition]
	data *genData.Data
}

func NewOrderLogisticsRepo(data *genData.Data) OrderLogisticsRepo {
	base := genRepo.NewBaseRepo[models.OrderLogistics, OrderLogisticsCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).OrderLogistics.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).OrderLogistics.ID
		},
		func(entity *models.OrderLogistics) int64 {
			return entity.ID
		},
		new(models.OrderLogistics),
		100,
	)
	return &orderLogisticsRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *orderLogisticsRepo) DeleteByOrderIds(ctx context.Context, orderIds []int64) error {
	q := r.data.Query(ctx).OrderLogistics
	_, err := q.WithContext(ctx).Where(q.OrderID.In(orderIds...)).Delete()
	return err
}
