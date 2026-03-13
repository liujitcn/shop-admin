package biz

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	_string "github.com/liujitcn/go-utils/string"
	_time "github.com/liujitcn/go-utils/time"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/kratos-kit/auth"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/api/gen/go/pay"
	"github.com/liujitcn/shop-admin/server/internal/configs"
	"github.com/liujitcn/shop-admin/server/internal/data"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type PayCase struct {
	tx                 genData.Transaction
	orderRepo          data.OrderRepo
	orderGoodsRepo     data.OrderGoodsRepo
	orderPaymentRepo   data.OrderPaymentRepo
	orderRefundRepo    data.OrderRefundRepo
	orderSchedulerCase *OrderSchedulerCase
	wxPayCase          *WxPayCase
}

// NewPayCase new a ShopPay use case.
func NewPayCase(
	tx genData.Transaction,
	orderCase data.OrderRepo,
	orderGoodsRepo data.OrderGoodsRepo,
	orderPaymentRepo data.OrderPaymentRepo,
	orderRefundRepo data.OrderRefundRepo,
	orderSchedulerCase *OrderSchedulerCase,
	wxPayCase *WxPayCase,
) *PayCase {
	return &PayCase{
		tx:                 tx,
		orderRepo:          orderCase,
		orderGoodsRepo:     orderGoodsRepo,
		orderPaymentRepo:   orderPaymentRepo,
		orderRefundRepo:    orderRefundRepo,
		orderSchedulerCase: orderSchedulerCase,
		wxPayCase:          wxPayCase,
	}
}

func (c *PayCase) JsapiPay(ctx context.Context, req *pay.PayRequest) (*pay.JsapiPayResponse, error) {
	authInfo, err := auth.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	var order *models.Order
	order, err = c.orderRepo.Find(ctx, &data.OrderCondition{
		Id:     req.GetOrderId(),
		UserId: authInfo.UserId,
	})
	if err != nil {
		return nil, err
	}
	if order.Status != int32(common.OrderStatus_CREATED) {
		return nil, fmt.Errorf("订单状态错误：【%s】", common.OrderStatus_name[order.Status])
	}

	var goods []*models.OrderGoods
	goods, err = c.orderGoodsRepo.FindAll(ctx, &data.OrderGoodsCondition{
		OrderId: order.ID,
	})
	if err != nil {
		return nil, err
	}

	goodsDetail := make([]jsapi.GoodsDetail, 0)
	for _, item := range goods {
		goodsId := fmt.Sprintf("%s_%s", strconv.FormatInt(item.GoodsID, 10), item.SkuCode)
		goodsDetail = append(goodsDetail, jsapi.GoodsDetail{
			MerchantGoodsId: &goodsId,
			GoodsName:       &item.Name,
			Quantity:        &item.Num,
			UnitPrice:       &item.Price,
		})
	}

	payTimeout := configs.ParsePayTimeout()
	createdAt := order.CreatedAt.Add(payTimeout)

	var description = "小程序支付"
	if len(goodsDetail) > 0 {
		description = trans.StringValue(goodsDetail[0].GoodsName)
	}

	return c.wxPayCase.JsapiPay(jsapi.PrepayRequest{
		Description: &description,
		OutTradeNo:  &order.OrderNo,
		TimeExpire:  &createdAt,
		Amount: &jsapi.Amount{
			Total: &order.PayMoney,
		},
		Payer: &jsapi.Payer{
			Openid: &authInfo.OpenId,
		},
		Detail: &jsapi.Detail{
			GoodsDetail: goodsDetail,
		},
	})
}

func (c *PayCase) H5Pay(ctx context.Context, req *pay.PayRequest) (*pay.H5PayResponse, error) {
	authInfo, err := auth.FromContext(ctx)
	if err != nil {
		log.Errorf("用户认证失败[%s]", err.Error())
		return nil, common.ErrorAccessForbidden("用户认证失败")
	}
	var order *models.Order
	order, err = c.orderRepo.Find(ctx, &data.OrderCondition{
		Id:     req.GetOrderId(),
		UserId: authInfo.UserId,
	})
	if err != nil {
		return nil, err
	}
	if order.Status != int32(common.OrderStatus_CREATED) {
		return nil, fmt.Errorf("订单状态错误：【%s】", common.OrderStatus_name[order.Status])
	}

	var goods []*models.OrderGoods
	goods, err = c.orderGoodsRepo.FindAll(ctx, &data.OrderGoodsCondition{
		OrderId: order.ID,
	})
	if err != nil {
		return nil, err
	}

	goodsDetail := make([]h5.GoodsDetail, 0)
	for _, item := range goods {
		goodsId := fmt.Sprintf("%s_%s", strconv.FormatInt(item.GoodsID, 10), item.SkuCode)
		goodsDetail = append(goodsDetail, h5.GoodsDetail{
			MerchantGoodsId: &goodsId,
			GoodsName:       &item.Name,
			Quantity:        &item.Num,
			UnitPrice:       &item.Price,
		})
	}
	payTimeout := configs.ParsePayTimeout()
	createdAt := order.CreatedAt.Add(payTimeout)

	var description = "H5支付"
	if len(goodsDetail) > 0 {
		description = trans.StringValue(goodsDetail[0].GoodsName)
	}
	return c.wxPayCase.H5Pay(h5.PrepayRequest{
		Description: trans.String(description),
		OutTradeNo:  trans.String(order.OrderNo),
		TimeExpire:  trans.Time(createdAt),
		Amount: &h5.Amount{
			Total: &order.PayMoney,
		},
		SceneInfo: &h5.SceneInfo{
			PayerClientIp: nil,
			DeviceId:      nil,
			StoreInfo:     nil,
			H5Info: &h5.H5Info{
				Type: trans.String("Wap"),
			},
		},
		Detail: &h5.Detail{
			GoodsDetail: goodsDetail,
		},
	})
}

func (c *PayCase) PayNotify(ctx context.Context, req *emptypb.Empty) error {
	request, err := c.wxPayCase.Notify(ctx)
	if err != nil {
		return err
	}
	resource := request.Resource
	if resource == nil {
		return errors.New("notify resource is nil")
	}

	log.Infof("PayNotify EventType=%s，Plaintext=%s", request.EventType, resource.Plaintext)
	// 判断通知类型
	if strings.HasPrefix(request.EventType, pay.ResourceType_TRANSACTION.String()) {
		// 转换
		var paymentResource pay.PaymentResource
		err = protojson.Unmarshal([]byte(resource.Plaintext), &paymentResource)
		if err != nil {
			return err
		}
		// 查询订单
		var order *models.Order
		order, err = c.orderRepo.FindByOrderNo(ctx, paymentResource.GetOutTradeNo())
		if err != nil {
			return err
		}
		// 查询支付信息
		var orderPayment *models.OrderPayment
		orderPayment, err = c.orderPaymentRepo.Find(ctx, &data.OrderPaymentCondition{
			OrderId: order.ID,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				orderPayment = &models.OrderPayment{}
			} else {
				return err
			}
		}
		successTime := _time.TimestamppbToTime(paymentResource.GetSuccessTime())
		if successTime == nil {
			successTime = trans.Time(time.Now())
		}
		orderPayment.OrderID = order.ID
		orderPayment.OrderNo = paymentResource.GetOutTradeNo()
		orderPayment.ThirdOrderNo = paymentResource.GetTransactionId()
		orderPayment.TradeType = paymentResource.GetTradeType().String()
		orderPayment.TradeState = paymentResource.GetTradeState().String()
		orderPayment.TradeStateDesc = paymentResource.GetTradeStateDesc()
		orderPayment.BankType = paymentResource.GetBankType()
		orderPayment.SuccessTime = trans.TimeValue(successTime)
		orderPayment.Payer = _string.ConvertAnyToJsonString(paymentResource.GetPayer())
		orderPayment.Amount = _string.ConvertAnyToJsonString(paymentResource.GetAmount())
		orderPayment.SceneInfo = _string.ConvertAnyToJsonString(paymentResource.GetSceneInfo())
		orderPayment.Status = 1

		return c.tx.Transaction(ctx, func(ctx context.Context) error {
			// 添加支付信息
			if orderPayment.ID == 0 {
				err = c.orderPaymentRepo.Create(ctx, orderPayment)
				if err != nil {
					return err
				}
			} else {
				err = c.orderPaymentRepo.UpdateByID(ctx, orderPayment)
				if err != nil {
					return err
				}
			}
			// 支付成功，修改订单状态
			if orderPayment.TradeState == pay.PaymentResource_SUCCESS.String() {
				err = c.orderRepo.UpdateByUserIdAndIds(ctx, order.UserID, []int64{order.ID}, &models.Order{
					Status: int32(common.OrderStatus_PAID),
				})
				if err != nil {
					return err
				}
				// 删除自动取消
				c.orderSchedulerCase.DeleteScheduled(order.ID)
			}
			return nil
		})
	} else if strings.HasPrefix(request.EventType, pay.ResourceType_REFUND.String()) {
		// 转换
		var refundResource pay.RefundResource
		err = protojson.Unmarshal([]byte(resource.Plaintext), &refundResource)
		if err != nil {
			return err
		}
		// 查询订单
		var order *models.Order
		order, err = c.orderRepo.FindByOrderNo(ctx, refundResource.GetOutTradeNo())
		if err != nil {
			return err
		}
		// 查询支付信息
		var orderRefund *models.OrderRefund
		orderRefund, err = c.orderRefundRepo.Find(ctx, &data.OrderRefundCondition{
			OrderId: order.ID,
		})
		successTime := _time.TimestamppbToTime(refundResource.GetSuccessTime())
		if successTime == nil {
			successTime = trans.Time(time.Now())
		}
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				orderRefund = &models.OrderRefund{
					OrderID:    order.ID,
					RefundNo:   refundResource.GetOutRefundNo(),
					CreateTime: time.Now(),
				}
			} else {
				return err
			}
		}
		orderRefund.OrderNo = refundResource.GetOutTradeNo()
		orderRefund.ThirdOrderNo = refundResource.GetTransactionId()
		orderRefund.ThirdRefundNo = refundResource.GetRefundId()
		orderRefund.UserReceivedAccount = refundResource.GetUserReceivedAccount()
		orderRefund.SuccessTime = trans.TimeValue(successTime)
		orderRefund.RefundState = refundResource.GetRefundStatus().String()
		orderRefund.Amount = _string.ConvertAnyToJsonString(refundResource.GetAmount())
		orderRefund.Status = 1

		return c.tx.Transaction(ctx, func(ctx context.Context) error {
			// 添加退款信息
			if orderRefund.ID == 0 {
				err = c.orderRefundRepo.Create(ctx, orderRefund)
				if err != nil {
					return err
				}
			} else {
				err = c.orderRefundRepo.UpdateByID(ctx, orderRefund)
				if err != nil {
					return err
				}
			}
			// 支付成功，修改订单状态
			if orderRefund.RefundState == pay.RefundResource_SUCCESS.String() {
				err = c.orderRepo.UpdateByUserIdAndIds(ctx, order.UserID, []int64{order.ID}, &models.Order{
					Status: int32(common.OrderStatus_REFUNDING),
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
	} else {
		return errors.New("notify event type err")
	}
}
