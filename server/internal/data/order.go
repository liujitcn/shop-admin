package data

import (
	"context"
	"errors"

	"time"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-admin/server/internal/data/dto"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"

	//"github.com/liujitcn/shop-admin/server/api/common"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type OrderCondition struct {
	Id             int64      `search:"type:eq;column:id"`
	UserId         int64      `search:"type:eq;column:user_id"`
	OrderNo        string     `search:"type:contains;column:order_no"`
	Status         int32      `search:"type:eq;column:status"`
	PayType        int32      `search:"type:eq;column:pay_type"`    // 支付方式，1为在线支付，2为货到付款
	PayChannel     int32      `search:"type:eq;column:pay_channel"` // 支付渠道：支付渠道，1支付宝、2微信--支付方式为在线支付时，传值，为货到付款时，不传值
	StartCreatedAt *time.Time `search:"type:gte;column:created_at"` // 创建开始时间
	EndCreatedAt   *time.Time `search:"type:lte;column:created_at"` // 创建结束时间
}

type OrderRepo interface {
	baseRepo.BaseRepo[models.Order, OrderCondition]
	UpdateByUserIdAndId(ctx context.Context, userID int64, orderInfo *models.Order) error
	UpdateByUserIdAndIds(ctx context.Context, userId int64, ids []int64, orderInfo *models.Order) error
	FindByOrderNo(ctx context.Context, orderNo string) (*models.Order, error)
	Sum(ctx context.Context, condition *OrderCondition) (int64, error)
	OrderSummary(ctx context.Context, timeType int32, condition *OrderCondition) ([]*dto.OrderSummary, error)
}

type orderRepo struct {
	baseRepo.BaseRepo[models.Order, OrderCondition]
	data *genData.Data
}

func NewOrderRepo(data *genData.Data) OrderRepo {
	base := baseRepo.NewBaseRepo[models.Order, OrderCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).Order.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).Order.ID
		},
		func(entity *models.Order) int64 {
			return entity.ID
		},
		new(models.Order),
	)
	return &orderRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *orderRepo) UpdateByUserIdAndId(ctx context.Context, userID int64, orderInfo *models.Order) error {
	if orderInfo.ID == 0 {
		return errors.New("orderInfo can not update without id")
	}
	q := r.data.Query(ctx).Order
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userID)).Updates(orderInfo)
	return err
}

func (r *orderRepo) UpdateByUserIdAndIds(ctx context.Context, userId int64, ids []int64, orderInfo *models.Order) error {
	if len(ids) == 0 {
		return errors.New("orderInfo can not update without id")
	}
	q := r.data.Query(ctx).Order
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userId), q.ID.In(ids...)).Updates(orderInfo)
	return err
}

func (r *orderRepo) FindByOrderNo(ctx context.Context, orderNo string) (*models.Order, error) {
	condition := OrderCondition{OrderNo: orderNo}
	return r.Find(ctx, &condition)
}
func (r *orderRepo) Sum(ctx context.Context, condition *OrderCondition) (int64, error) {
	m := r.data.Query(ctx).Order
	q, err := baseRepo.BuildDao(new(r.data.Query(ctx).Order.WithContext(ctx).DO), new(models.Order), condition)
	if err != nil {
		return 0, err
	}
	var results struct{ Total int64 }
	q.Select(m.PayMoney.Sum().As("total"))
	err = q.Scan(&results)
	if err != nil {
		return 0, err
	}
	return results.Total, nil
}

func (r *orderRepo) OrderSummary(ctx context.Context, timeType int32, condition *OrderCondition) ([]*dto.OrderSummary, error) {
	m := r.data.Query(ctx).Order
	q := m.WithContext(ctx)

	if condition.StartCreatedAt != nil {
		q = q.Where(m.CreatedAt.Gte(*condition.StartCreatedAt))
	}
	if condition.EndCreatedAt != nil {
		q = q.Where(m.CreatedAt.Lt(*condition.EndCreatedAt))
	}
	if condition.Status > 0 {
		q = q.Where(m.Status.Eq(condition.Status))
	}
	results := make([]*dto.OrderSummary, 0)
	switch timeType {
	case 2:
		q.Select(m.CreatedAt.Month().As("key"),
			m.ID.Count().As("order_count"),
			m.PayMoney.Sum().As("sale_amount")).Group(m.CreatedAt.Month()).Order(m.CreatedAt.Month())
	case 1:
		q.Select(m.CreatedAt.DayOfWeek().As("key"),
			m.ID.Count().As("order_count"),
			m.PayMoney.Sum().As("sale_amount")).Group(m.CreatedAt.DayOfWeek()).Order(m.CreatedAt.DayOfWeek())
	default:
		q.Select(m.CreatedAt.DayOfMonth().As("key"),
			m.ID.Count().As("order_count"),
			m.PayMoney.Sum().As("sale_amount")).Group(m.CreatedAt.DayOfMonth()).Order(m.CreatedAt.DayOfMonth())
	}
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
