package biz

import (
	"context"

	"github.com/liujitcn/go-utils/str"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/internal/data"
)

type OrderGoodsCase struct {
	data.OrderGoodsRepo
}

// NewOrderGoodsCase new a OrderGoods use case.
func NewOrderGoodsCase(orderGoodsRepo data.OrderGoodsRepo,
) *OrderGoodsCase {
	return &OrderGoodsCase{
		OrderGoodsRepo: orderGoodsRepo,
	}
}

func (c *OrderGoodsCase) GetFromByOrderId(ctx context.Context, orderId int64) ([]*admin.OrderGoods, error) {
	orderGoods, err := c.FindAll(ctx, &data.OrderGoodsCondition{
		OrderId: orderId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*admin.OrderGoods, 0)
	for _, item := range orderGoods {
		list = append(list, &admin.OrderGoods{
			GoodsId:       item.GoodsID,
			SkuCode:       item.SkuCode,
			Picture:       item.Picture,
			Name:          item.Name,
			Num:           item.Num,
			SpecItem:      str.ConvertJsonStringToStringArray(item.SpecItem),
			Price:         item.Price,
			PayPrice:      item.PayPrice,
			TotalPrice:    item.TotalPrice,
			TotalPayPrice: item.TotalPayPrice,
		})
	}

	return list, nil
}
