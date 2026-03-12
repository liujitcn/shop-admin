package biz

import (
	"context"
	"encoding/json"
	"errors"

	_time "github.com/liujitcn/go-utils/time"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"gorm.io/gorm"
)

type OrderPaymentCase struct {
	data.OrderPaymentRepo
}

// NewOrderPaymentCase new a OrderPayment use case.
func NewOrderPaymentCase(orderPaymentRepo data.OrderPaymentRepo,
) *OrderPaymentCase {
	return &OrderPaymentCase{
		OrderPaymentRepo: orderPaymentRepo,
	}
}

func (c *OrderPaymentCase) GetFromByOrderId(ctx context.Context, orderId int64) (*admin.OrderPayment, error) {
	orderPayment, err := c.Find(ctx, &data.OrderPaymentCondition{
		OrderId: orderId,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &admin.OrderPayment{}, nil
		}
		return nil, err
	}
	var payer admin.OrderPayment_Payer
	_ = json.Unmarshal([]byte(orderPayment.Payer), &payer)
	var amount admin.OrderPayment_Amount
	_ = json.Unmarshal([]byte(orderPayment.Amount), &amount)
	var sceneInfo admin.OrderPayment_SceneInfo
	_ = json.Unmarshal([]byte(orderPayment.SceneInfo), &sceneInfo)
	return &admin.OrderPayment{
		OrderNo:        orderPayment.OrderNo,
		ThirdOrderNo:   orderPayment.ThirdOrderNo,
		TradeType:      orderPayment.TradeType,
		TradeState:     orderPayment.TradeState,
		TradeStateDesc: orderPayment.TradeStateDesc,
		BankType:       orderPayment.BankType,
		SuccessTime:    _time.TimeToTimeString(orderPayment.SuccessTime),
		Payer:          &payer,
		Amount:         &amount,
		SceneInfo:      &sceneInfo,
		Status:         common.OrderBillStatus(orderPayment.Status),
	}, nil
}
