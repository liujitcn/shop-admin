package biz

import (
	"context"
	"encoding/json"
	"errors"

	_time "github.com/liujitcn/go-utils/time"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"gorm.io/gorm"
)

type OrderLogisticsCase struct {
	data.OrderLogisticsRepo
}

// NewOrderLogisticsCase new a OrderLogistics use case.
func NewOrderLogisticsCase(orderLogisticsRepo data.OrderLogisticsRepo,
) *OrderLogisticsCase {
	return &OrderLogisticsCase{
		OrderLogisticsRepo: orderLogisticsRepo,
	}
}

func (c *OrderLogisticsCase) GetFromByOrderId(ctx context.Context, orderId int64) (*admin.OrderLogistics, error) {
	orderLogistics, err := c.Find(ctx, &data.OrderLogisticsCondition{
		OrderId: orderId,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &admin.OrderLogistics{}, nil
		}
		return nil, err
	}
	detail := make([]*admin.OrderLogistics_Detail, 0)
	_ = json.Unmarshal([]byte(orderLogistics.Detail), &detail)
	return &admin.OrderLogistics{
		Name:      orderLogistics.Name,
		No:        orderLogistics.No,
		Contact:   orderLogistics.Contact,
		Detail:    detail,
		CreatedAt: _time.TimeToTimeString(orderLogistics.CreatedAt),
	}, nil
}
