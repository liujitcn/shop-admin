package biz

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	nethttp "net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/conf"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	"github.com/liujitcn/shop-admin/server/internal/service/pay/bill"
	wxPayCore "github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type WxPayCase struct {
	wxPay         *conf.WxPay
	mchPrivateKey *rsa.PrivateKey
	ctx           context.Context
	client        *wxPayCore.Client
}

// NewWxPayCase 微信支付
func NewWxPayCase(
	wxPay *conf.WxPay,
) (*WxPayCase, error) {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(wxPay.GetMchCertPath())
	if err != nil {
		return nil, err
	}
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []wxPayCore.ClientOption{
		option.WithWechatPayAutoAuthCipher(wxPay.GetMchId(), wxPay.GetMchCertSn(), mchPrivateKey, wxPay.GetMchAPIv3Key()),
	}
	ctx := context.Background()
	var client *wxPayCore.Client
	client, err = wxPayCore.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &WxPayCase{
		wxPay:         wxPay,
		mchPrivateKey: mchPrivateKey,
		ctx:           ctx,
		client:        client,
	}, nil
}

func (c *WxPayCase) Refund(req refunddomestic.CreateRequest) (*refunddomestic.Refund, error) {
	// 拼接公共参数
	req.NotifyUrl = trans.String(c.wxPay.GetNotifyUrl() + _const.NotifyUrl)

	svc := refunddomestic.RefundsApiService{Client: c.client}
	resp, result, err := svc.Create(c.ctx, req)
	if err != nil {
		log.Errorf("支付失败[%s]", err.Error())
		return nil, errors.New("支付失败")
	}
	if result.Response.StatusCode != nethttp.StatusOK {
		log.Errorf("支付失败[%s]", result.Response.Status)
		return nil, errors.New("支付失败")
	}

	return resp, err
}

func (c *WxPayCase) QueryByOutRefundNo(req refunddomestic.QueryByOutRefundNoRequest) (*refunddomestic.Refund, error) {
	// 拼接公共参数
	svc := refunddomestic.RefundsApiService{Client: c.client}
	resp, result, err := svc.QueryByOutRefundNo(c.ctx, req)
	if err != nil {
		log.Errorf("查询退款失败[%s]", err.Error())
		return nil, errors.New("查询退款失败")
	}
	if result.Response.StatusCode != nethttp.StatusOK {
		log.Errorf("查询退款失败[%s]", result.Response.Status)
		return nil, errors.New("查询退款失败")
	}

	return resp, err
}

func (c *WxPayCase) QueryOrderByOutTradeNo(req jsapi.QueryOrderByOutTradeNoRequest) (*payments.Transaction, error) {
	req.Mchid = trans.String(c.wxPay.GetMchId())
	svc := jsapi.JsapiApiService{Client: c.client}
	resp, result, err := svc.QueryOrderByOutTradeNo(c.ctx, req)
	if err != nil {
		log.Errorf("查询支付失败[%s]", err.Error())
		return nil, errors.New("查询支付失败")
	}
	if result.Response.StatusCode != nethttp.StatusOK {
		log.Errorf("查询支付失败[%s]", result.Response.Status)
		return nil, errors.New("查询支付失败")
	}

	return resp, err
}

func (c *WxPayCase) TradeBill(req bill.TradeBillRequest) (*bill.TradeBillResponse, error) {
	svc := bill.BillApiService{Client: c.client}
	resp, result, err := svc.TradeBill(c.ctx, req)
	if err != nil {
		log.Errorf("申请交易账单失败[%s]", err.Error())
		return nil, errors.New("申请交易账单失败")
	}
	if result.Response.StatusCode != nethttp.StatusOK {
		log.Errorf("申请交易账单失败[%s]", result.Response.Status)
		return nil, errors.New("申请交易账单失败")
	}
	return resp, err
}

func (c *WxPayCase) DownloadBill(url string) ([]byte, error) {
	svc := bill.BillApiService{Client: c.client}
	// 下载账单
	return svc.DownloadBill(c.ctx, url)
}

func (c *WxPayCase) Notify(ctx context.Context) (*notify.Request, error) {
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err := downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, c.mchPrivateKey, c.wxPay.GetMchCertSn(), c.wxPay.GetMchId(), c.wxPay.GetMchAPIv3Key())
	if err != nil {
		return nil, err
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(c.wxPay.GetMchId())

	// 3. 使用证书访问器初始化 `notify.Handler`
	var handler *notify.Handler
	handler, err = notify.NewRSANotifyHandler(c.wxPay.GetMchAPIv3Key(), verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		return nil, err
	}
	var httpReq *nethttp.Request
	if info, ok := transport.FromServerContext(ctx); ok {
		if htr, htrOk := info.(*http.Transport); htrOk {
			httpReq = htr.Request()
		}
	}
	if httpReq == nil {
		return nil, errors.New("transport convert nethttp request failed")
	}
	var req *notify.Request
	req, err = handler.ParseNotifyRequest(ctx, httpReq, certificateVisitor)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *WxPayCase) generatePaySign(timeStamp, nonceStr, packageStr string) string {

	var signBuilder strings.Builder
	signBuilder.WriteString(c.wxPay.GetAppid() + "\n")
	signBuilder.WriteString(timeStamp + "\n")
	signBuilder.WriteString(nonceStr + "\n")
	signBuilder.WriteString(packageStr + "\n")
	signString := signBuilder.String()

	hashed := sha256.Sum256([]byte(signString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, c.mchPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}
