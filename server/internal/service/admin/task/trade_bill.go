package task

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/liujitcn/go-sdk/oss"
	"github.com/liujitcn/go-utils"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-admin/server/internal/service/pay/bill"
	"github.com/liujitcn/shop-admin/server/internal/service/pay/biz"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gorm"
)

// TradeBill 申请交易账单
type TradeBill struct {
	oss              oss.OSS
	tx               genData.Transaction
	wxPayCase        *biz.WxPayCase
	payBillCase      *biz.PayBillCase
	orderPaymentRepo data.OrderPaymentRepo
	orderRefundRepo  data.OrderRefundRepo
	ctx              context.Context
}

func NewTradeBill(
	oss oss.OSS,
	tx genData.Transaction,
	wxPayCase *biz.WxPayCase,
	payBillCase *biz.PayBillCase,
	orderPaymentRepo data.OrderPaymentRepo,
	orderRefundRepo data.OrderRefundRepo,
) *TradeBill {
	return &TradeBill{
		oss:              oss,
		tx:               tx,
		wxPayCase:        wxPayCase,
		payBillCase:      payBillCase,
		orderPaymentRepo: orderPaymentRepo,
		orderRefundRepo:  orderRefundRepo,
		ctx:              context.Background(),
	}
}

func (t *TradeBill) Exec(args map[string]string) ([]string, error) {
	log.Infof("Job TradeBill Exec %+v", args)
	var billDate string
	v, ok := args["billDate"]
	if ok {
		billDate = v
	} else {
		billDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	ret := make([]string, 0)
	payment, err1 := t.payment(billDate, bill.BILL_TYPE_SUCCESS)
	if err1 != nil {
		return nil, err1
	}
	ret = append(ret, payment...)
	refund, err2 := t.refund(billDate, bill.BILL_TYPE_REFUND)
	if err2 != nil {
		return nil, err2
	}
	ret = append(ret, refund...)
	// 下载账单
	return ret, nil
}

func (t *TradeBill) payment(billDate, billType string) ([]string, error) {
	ret := make([]string, 0)
	payBill, err := t.downloadBill(billDate, billType)
	if err != nil {
		ret = append(ret, err.Error())
		return ret, err
	}
	// 查询全部支付订单
	paymentList := make([]*models.OrderPayment, 0)
	paymentList, err = t.orderPaymentRepo.FindAll(t.ctx, &data.OrderPaymentCondition{
		SuccessTime: billDate,
	})
	if err != nil {
		ret = append(ret, err.Error())
		return ret, err
	}
	// 转换map
	paymentMap := make(map[string]*models.OrderPayment)
	for _, payment := range paymentList {
		paymentMap[fmt.Sprintf("%s_%s", payment.OrderNo, payment.ThirdOrderNo)] = payment
	}

	// 获取对账单内容
	var fileByte []byte
	fileByte, err = t.oss.GetFileByte(payBill.FilePath)
	if err != nil {
		ret = append(ret, err.Error())
		return ret, err
	}
	err = t.checkHash(fileByte, payBill.HashValue)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(fileByte))
	reader.Comma = ','       // 设置分隔符
	reader.LazyQuotes = true // 允许非常规引号

	// 跳过标题行
	_, _ = reader.Read()

	for {
		var record []string
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		// 去除每个字段的反引号
		for i := range record {
			record[i] = strings.Trim(record[i], "`")
		}
		switch len(record) {
		case 20:
			// 计算金额
			amount := utils.ConvertYuanToFen(record[12])
			payBill.ThirdTotalCount += 1
			payBill.ThirdTotalAmount += amount
			// 交易记录
			key := fmt.Sprintf("%s_%s", record[6], record[5])
			// 记录在数据库不存在，暂时记录日期，后续在做处理
			if v, ok := paymentMap[key]; ok {
				var orderPaymentAmount admin.OrderPayment_Amount
				_ = json.Unmarshal([]byte(v.Amount), &orderPaymentAmount)
				// 支付金额和状态一致
				if v.TradeState == record[9] && orderPaymentAmount.GetPayerTotal() == amount {
					v.Status = 2
				} else {
					v.Status = 3
				}
			} else {
				ret = append(ret, fmt.Sprintf("%+v", record))
			}
		default:
			continue
		}
	}
	err = t.tx.Transaction(t.ctx, func(ctx context.Context) error {
		for _, v := range paymentMap {
			payBill.TotalCount += 1
			var orderPaymentAmount admin.OrderPayment_Amount
			_ = json.Unmarshal([]byte(v.Amount), &orderPaymentAmount)
			payBill.TotalAmount += orderPaymentAmount.GetPayerTotal()
			err = t.orderPaymentRepo.UpdateByID(ctx, v)
			if err != nil {
				return err
			}
		}

		if payBill.TotalCount == payBill.ThirdTotalCount && payBill.TotalAmount == payBill.ThirdTotalAmount {
			payBill.Status = 1
		} else {
			payBill.Status = 2
		}
		return t.payBillCase.UpdateByID(ctx, payBill)
	})
	if err != nil {
		ret = append(ret, err.Error())
		return ret, err
	}

	return ret, nil
}

func (t *TradeBill) refund(billDate, billType string) ([]string, error) {
	ret := make([]string, 0)
	payBill, err := t.downloadBill(billDate, billType)
	if err != nil {
		ret = append(ret, err.Error())
		return ret, err
	}
	// 查询全部退款订单
	refundList := make([]*models.OrderRefund, 0)
	refundList, err = t.orderRefundRepo.FindAll(t.ctx, &data.OrderRefundCondition{
		SuccessTime: billDate,
	})
	if err != nil {
		ret = append(ret, err.Error())
		return ret, err
	}
	// 转换map
	refundMap := make(map[string]*models.OrderRefund)
	for _, refund := range refundList {
		refundMap[fmt.Sprintf("%s_%s_%s_%s", refund.OrderNo, refund.ThirdOrderNo, refund.RefundNo, refund.ThirdRefundNo)] = refund
	}

	// 获取对账单内容
	var fileByte []byte
	fileByte, err = t.oss.GetFileByte(payBill.FilePath)
	if err != nil {
		ret = append(ret, err.Error())
		return ret, err
	}
	err = t.checkHash(fileByte, payBill.HashValue)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(fileByte))
	reader.Comma = ','       // 设置分隔符
	reader.LazyQuotes = true // 允许非常规引号

	// 跳过标题行
	_, _ = reader.Read()

	for {
		var record []string
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		// 去除每个字段的反引号
		for i := range record {
			record[i] = strings.Trim(record[i], "`")
		}
		switch len(record) {
		case 29:
			// 计算金额
			amount := utils.ConvertYuanToFen(record[18])
			payBill.ThirdTotalCount += 1
			payBill.ThirdTotalAmount += amount
			// 交易记录
			key := fmt.Sprintf("%s_%s_%s_%s", record[6], record[5], record[17], record[16])
			// 记录在数据库不存在，暂时记录日期，后续在做处理
			if v, ok := refundMap[key]; ok {
				var orderRefundAmount admin.OrderRefund_Amount
				_ = json.Unmarshal([]byte(v.Amount), &orderRefundAmount)
				// 支付金额和状态一致
				if v.RefundState == record[21] && orderRefundAmount.GetPayerRefund() == amount {
					v.Status = 2
				} else {
					v.Status = 3
				}
			} else {
				ret = append(ret, fmt.Sprintf("%+v", record))
			}
		default:
			continue
		}
	}
	err = t.tx.Transaction(t.ctx, func(ctx context.Context) error {
		for _, v := range refundMap {
			payBill.TotalCount += 1
			var orderRefundAmount admin.OrderRefund_Amount
			_ = json.Unmarshal([]byte(v.Amount), &orderRefundAmount)
			payBill.TotalAmount += orderRefundAmount.GetPayerRefund()
			err = t.orderRefundRepo.UpdateByID(ctx, v)
			if err != nil {
				return err
			}
		}

		if payBill.TotalCount == payBill.ThirdTotalCount && payBill.TotalAmount == payBill.ThirdTotalAmount {
			payBill.Status = 1
		} else {
			payBill.Status = 2
		}
		return t.payBillCase.UpdateByID(ctx, payBill)
	})
	if err != nil {
		ret = append(ret, err.Error())
		return ret, err
	}
	return ret, nil
}

// 下载对账单
func (t *TradeBill) downloadBill(billDate, billType string) (*models.PayBill, error) {
	// 获取当前定时账单日期
	first, err := t.payBillCase.Find(t.ctx, &data.PayBillCondition{
		BillDate: billDate,
		BillType: billType,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || first == nil {
		// 申请账单
		var tradeBill *bill.TradeBillResponse
		tradeBill, err = t.wxPayCase.TradeBill(bill.TradeBillRequest{
			BillDate: &billDate,
			BillType: &billType,
		})
		if err != nil {
			return nil, err
		}

		// 下载账单
		var billByte []byte
		billByte, err = t.wxPayCase.DownloadBill(trans.StringValue(tradeBill.DownloadUrl))
		if err != nil {
			return nil, err
		}

		// 校验Hash
		hashValue := trans.StringValue(tradeBill.HashValue)
		err = t.checkHash(billByte, hashValue)
		if err != nil {
			return nil, err
		}

		var path string
		path, err = t.oss.UploadByByte(fmt.Sprintf("%s.csv", billType), fmt.Sprintf("bill/file/%s", strings.ReplaceAll(billDate, "-", "/")), billByte)
		if err != nil {
			return nil, err
		}

		first = &models.PayBill{
			BillDate:  billDate,
			BillType:  billType,
			FilePath:  path,
			HashType:  trans.StringValue(tradeBill.HashType),
			HashValue: hashValue,
		}
		err = t.payBillCase.Create(t.ctx, first)
		if err != nil {
			return nil, err
		}
	} else {
		// 重新计算
		first.TotalCount = 0
		first.TotalAmount = 0
		first.ThirdTotalCount = 0
		first.ThirdTotalAmount = 0
		first.Status = 0
	}
	return first, nil
}

func (t *TradeBill) checkHash(fileBytes []byte, hashValue string) error {
	hash := sha1.New()
	hash.Write(fileBytes)
	hashSum := hash.Sum(nil) // 返回 [20]byte 的切片
	// 将哈希值转换为十六进制字符串（与常见工具格式一致）
	hashHex := fmt.Sprintf("%x", hashSum)
	if hashHex != hashValue {
		return errors.New("hash value error")
	}
	return nil
}
