package biz

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/liujitcn/go-utils/id"
	_string "github.com/liujitcn/go-utils/string"
	_time "github.com/liujitcn/go-utils/time"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/api/gen/go/pay"
	"github.com/liujitcn/shop-admin/server/internal/configs"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	"github.com/liujitcn/shop-admin/server/internal/data"
	payBiz "github.com/liujitcn/shop-admin/server/internal/service/pay/biz"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"gorm.io/gorm"
)

type OrderCase struct {
	tx genData.Transaction
	data.OrderRepo
	orderAddressCase   *OrderAddressCase
	orderCancelCase    *OrderCancelCase
	orderGoodsCase     *OrderGoodsCase
	orderLogisticsCase *OrderLogisticsCase
	orderPaymentCase   *OrderPaymentCase
	orderRefundCase    *OrderRefundCase
	baseUserCase       *BaseUserCase
	baseDictItemCase   *BaseDictItemCase
	wxPayCase          *payBiz.WxPayCase
}

// NewOrderCase new a Order use case.
func NewOrderCase(
	tx genData.Transaction,
	orderAddressCase *OrderAddressCase,
	OrderRepo data.OrderRepo,
	orderCancelCase *OrderCancelCase,
	orderGoodsCase *OrderGoodsCase,
	orderLogisticsCase *OrderLogisticsCase,
	orderPaymentCase *OrderPaymentCase,
	orderRefundCase *OrderRefundCase,
	baseUserCase *BaseUserCase,
	baseDictItemCase *BaseDictItemCase,
	wxPayCase *payBiz.WxPayCase,
) *OrderCase {
	return &OrderCase{
		tx:                 tx,
		OrderRepo:          OrderRepo,
		orderAddressCase:   orderAddressCase,
		orderCancelCase:    orderCancelCase,
		orderGoodsCase:     orderGoodsCase,
		orderLogisticsCase: orderLogisticsCase,
		orderPaymentCase:   orderPaymentCase,
		orderRefundCase:    orderRefundCase,
		baseUserCase:       baseUserCase,
		baseDictItemCase:   baseDictItemCase,
		wxPayCase:          wxPayCase,
	}
}
func (c *OrderCase) GetFromID(ctx context.Context, id int64) (*admin.OrderResponse, error) {
	var err error
	var order *models.Order
	order, err = c.Find(ctx, &data.OrderCondition{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	_order := c.ConvertToProto(order)

	// 查询用户
	var baseUser *models.BaseUser
	baseUser, err = c.baseUserCase.GetFromID(ctx, order.UserID)
	if err != nil {
		return nil, err
	}
	_order.NickName = baseUser.NickName

	// 支付超时时间
	payTimeout := configs.ParsePayTimeout()
	createdAt := order.CreatedAt.Add(payTimeout)
	nowTime := time.Now()
	countdown := float32(createdAt.Sub(nowTime).Seconds())

	res := &admin.OrderResponse{
		Order:     _order,
		Countdown: countdown,
	}
	// 地址
	res.Address, err = c.orderAddressCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	// 取消
	res.Cancel, err = c.orderCancelCase.GetFromByOrderId(ctx, order.ID)

	// 商品
	res.Goods, err = c.orderGoodsCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	// 发货
	res.Logistics, err = c.orderLogisticsCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	// 在线支付
	res.Payment, err = c.orderPaymentCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	// 退款
	res.Refund, err = c.orderRefundCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *OrderCase) Page(ctx context.Context, req *admin.PageOrderRequest) (*admin.PageOrderResponse, error) {
	var err error
	var startCreatedAt, endCreatedAt *time.Time
	createdAt := req.GetCreatedAt()
	if len(createdAt) == 2 {
		startCreatedAt = _time.StringTimeToTime(createdAt[0])
		endCreatedAt = _time.StringTimeToTime(createdAt[1])
		if endCreatedAt != nil {
			t := endCreatedAt.AddDate(0, 0, 1)
			endCreatedAt = &t
		}
	}
	condition := &data.OrderCondition{
		UserId:         req.GetUserId(),
		OrderNo:        req.GetOrderNo(),
		Status:         int32(req.GetStatus()),
		PayType:        int32(req.GetPayType()),
		PayChannel:     int32(req.GetPayChannel()),
		StartCreatedAt: startCreatedAt,
		EndCreatedAt:   endCreatedAt,
	}
	var page []*models.Order
	var count int64
	page, count, err = c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.Order, 0)
	for _, item := range page {
		Order := c.ConvertToProto(item)
		list = append(list, Order)
	}

	return &admin.PageOrderResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *OrderCase) GetOrderRefund(ctx context.Context, id int64) (*admin.OrderRefundResponse, error) {
	order, err := c.Find(ctx, &data.OrderCondition{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	res := &admin.OrderRefundResponse{}
	// 支付
	res.Payment, err = c.orderPaymentCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	// 退款
	res.Refund, err = c.orderRefundCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *OrderCase) RefundOrder(ctx context.Context, req *admin.RefundOrderRequest) error {
	order, err := c.Find(ctx, &data.OrderCondition{
		Id: req.GetOrderId(),
	})
	if err != nil {
		return err
	}
	if !(order.Status == int32(common.OrderStatus_SHIPPED) || order.Status == int32(common.OrderStatus_RECEIVED) || order.Status == int32(common.OrderStatus_REFUNDING)) {
		return fmt.Errorf("订单状态错误：【%s】", common.OrderStatus_name[order.Status])
	}

	orderRefund := &models.OrderRefund{
		OrderID:  req.GetOrderId(),
		RefundNo: strconv.FormatInt(id.GenSnowflakeID(), 10),
		Reason:   int32(req.GetReason()),
	}
	// 在线支付
	if common.OrderPayType(order.PayType) == common.OrderPayType_ONLINE_PAY {
		// 查询退款原因
		reason := strconv.Itoa(int(orderRefund.Reason))
		var label string
		label, err = c.baseDictItemCase.FindLabelByCodeAndValue(ctx, _const.BaseDictCode_OrderRefundReason, reason)
		if err == nil {
			reason = label
		}
		// 微信支付
		if common.OrderPayChannel(order.PayChannel) == common.OrderPayChannel_WX_PAY {
			var refund *refunddomestic.Refund
			refund, err = c.wxPayCase.Refund(refunddomestic.CreateRequest{
				OutTradeNo:  trans.String(order.OrderNo),
				OutRefundNo: trans.String(orderRefund.RefundNo),
				Reason:      trans.String(reason),
				Amount: &refunddomestic.AmountReq{
					Total:    trans.Int64(order.PayMoney),
					Refund:   trans.Int64(req.GetRefundMoney()),
					Currency: trans.String("CNY"),
				},
			})
			if err != nil {
				return err
			}
			orderRefund.OrderNo = trans.StringValue(refund.OutTradeNo)
			orderRefund.ThirdOrderNo = trans.StringValue(refund.TransactionId)
			orderRefund.ThirdRefundNo = trans.StringValue(refund.RefundId)
			orderRefund.Channel = string(*refund.Channel.Ptr())
			orderRefund.UserReceivedAccount = trans.StringValue(refund.UserReceivedAccount)
			orderRefund.CreateTime = trans.TimeValue(refund.CreateTime)
			orderRefund.SuccessTime = trans.TimeValue(refund.SuccessTime)
			orderRefund.RefundState = string(*refund.Status.Ptr())
			orderRefund.FundsAccount = string(*refund.FundsAccount)
			orderRefund.Amount = _string.ConvertAnyToJsonString(refund.Amount)
			orderRefund.Status = 1
		}
	} else {
		t := time.Now()
		orderRefund.CreateTime = t
		orderRefund.SuccessTime = t
		orderRefund.Amount = "{}"
	}
	ids := []int64{req.GetOrderId()}
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		// 增加退款信息
		err = c.orderRefundCase.Create(ctx, orderRefund)
		if err != nil {
			return err
		}
		return c.UpdateByUserIdAndIds(ctx, order.UserID, ids, &models.Order{
			Status: int32(common.OrderStatus_REFUNDING),
		})
	})
}

func (c *OrderCase) GetOrderShipped(ctx context.Context, id int64) (*admin.OrderShippedResponse, error) {
	order, err := c.Find(ctx, &data.OrderCondition{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	// 商品
	var goods []*admin.OrderGoods
	goods, err = c.orderGoodsCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	// 地址
	var address *admin.OrderAddress
	address, err = c.orderAddressCase.GetFromByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	res := &admin.OrderShippedResponse{
		Address: address,
		Goods:   goods,
	}
	// 判断订单状态
	switch order.Status {
	case int32(common.OrderStatus_SHIPPED), int32(common.OrderStatus_RECEIVED):
		// 发货
		res.Logistics, err = c.orderLogisticsCase.GetFromByOrderId(ctx, order.ID)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (c *OrderCase) ShippedOrder(ctx context.Context, req *admin.ShippedOrderRequest) error {
	order, err := c.Find(ctx, &data.OrderCondition{
		Id: req.GetOrderId(),
	})
	if err != nil {
		return err
	}
	if order.Status != int32(common.OrderStatus_PAID) {
		return fmt.Errorf("订单状态错误：【%s】", common.OrderStatus_name[order.Status])
	}
	// 查询商品是否支付
	// 在线支付
	if common.OrderPayType(order.PayType) == common.OrderPayType_ONLINE_PAY {
		// 微信支付
		if common.OrderPayChannel(order.PayChannel) == common.OrderPayChannel_WX_PAY {
			var transaction *payments.Transaction
			transaction, err = c.wxPayCase.QueryOrderByOutTradeNo(jsapi.QueryOrderByOutTradeNoRequest{
				OutTradeNo: trans.String(order.OrderNo),
			})
			if err != nil {
				return err
			}
			tradeState := trans.StringValue(transaction.TradeState)
			if tradeState != pay.PaymentResource_SUCCESS.String() {
				return fmt.Errorf("订单状态错误：【%s】", tradeState)
			}
			// 查询支付信息
			var orderPayment *models.OrderPayment
			orderPayment, err = c.orderPaymentCase.Find(ctx, &data.OrderPaymentCondition{
				OrderId: order.ID,
			})
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					orderPayment = &models.OrderPayment{}
				} else {
					return err
				}
			}
			successTime := _time.StringDateToTime(transaction.SuccessTime)
			if successTime == nil {
				successTime = trans.Time(time.Now())
			}
			orderPayment.OrderID = order.ID
			orderPayment.OrderNo = trans.StringValue(transaction.OutTradeNo)
			orderPayment.ThirdOrderNo = trans.StringValue(transaction.TransactionId)
			orderPayment.TradeType = trans.StringValue(transaction.TradeType)
			orderPayment.TradeState = tradeState
			orderPayment.TradeStateDesc = trans.StringValue(transaction.TradeStateDesc)
			orderPayment.BankType = trans.StringValue(transaction.BankType)
			orderPayment.SuccessTime = trans.TimeValue(successTime)
			orderPayment.Payer = _string.ConvertAnyToJsonString(transaction.Payer)
			orderPayment.Amount = _string.ConvertAnyToJsonString(transaction.Amount)
			// 添加支付信息
			if orderPayment.ID == 0 {
				err = c.orderPaymentCase.Create(ctx, orderPayment)
				if err != nil {
					return err
				}
			} else {
				err = c.orderPaymentCase.UpdateByID(ctx, orderPayment)
				if err != nil {
					return err
				}
			}
		}
	}
	ids := []int64{req.GetOrderId()}
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		// 增加退款信息
		err = c.orderLogisticsCase.Create(ctx, &models.OrderLogistics{
			OrderID: order.ID,
			Name:    req.GetName(),
			No:      req.GetNo(),
			Contact: req.GetContact(),
			Detail:  "[]",
		})
		if err != nil {
			return err
		}
		return c.UpdateByUserIdAndIds(ctx, order.UserID, ids, &models.Order{
			Status: int32(common.OrderStatus_SHIPPED),
		})
	})
}

func (c *OrderCase) ConvertToProto(item *models.Order) *admin.Order {
	res := &admin.Order{
		Id:           item.ID,
		OrderNo:      item.OrderNo,
		UserId:       item.UserID,
		PayMoney:     item.PayMoney,
		TotalMoney:   item.TotalMoney,
		PostFee:      item.PostFee,
		GoodsNum:     item.GoodsNum,
		PayType:      common.OrderPayType(item.PayType),
		PayChannel:   common.OrderPayChannel(item.PayChannel),
		DeliveryTime: common.OrderDeliveryTime(item.DeliveryTime),
		Status:       common.OrderStatus(item.Status),
		Remark:       item.Remark,
		CreatedAt:    _time.TimeToTimeString(item.CreatedAt),
		UpdatedAt:    _time.TimeToTimeString(item.UpdatedAt),
	}
	return res
}
