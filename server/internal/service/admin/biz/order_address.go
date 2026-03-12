package biz

import (
	"context"

	_string "github.com/liujitcn/go-utils/string"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/internal/data"
)

type OrderAddressCase struct {
	data.OrderAddressRepo
}

// NewOrderAddressCase new a OrderAddress use case.
func NewOrderAddressCase(orderAddressRepo data.OrderAddressRepo,
) *OrderAddressCase {
	return &OrderAddressCase{
		OrderAddressRepo: orderAddressRepo,
	}
}

func (c *OrderAddressCase) GetFromByOrderId(ctx context.Context, orderId int64) (*admin.OrderAddress, error) {
	orderAddress, err := c.Find(ctx, &data.OrderAddressCondition{
		OrderId: orderId,
	})
	if err != nil {
		return nil, err
	}
	return &admin.OrderAddress{
		Receiver: orderAddress.Receiver,
		Contact:  orderAddress.Contact,
		Address:  _string.ConvertJsonStringToStringArray(orderAddress.Address),
		Detail:   orderAddress.Detail,
	}, nil
}
