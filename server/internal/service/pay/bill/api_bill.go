package bill

import (
	"context"
	"fmt"
	"io"
	nethttp "net/http"
	neturl "net/url"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/validators"
	"github.com/wechatpay-apiv3/wechatpay-go/core/consts"
	"github.com/wechatpay-apiv3/wechatpay-go/services"
)

type BillApiService services.Service

// TradeBill 申请交易账单
func (a *BillApiService) TradeBill(ctx context.Context, req TradeBillRequest) (resp *TradeBillResponse, result *core.APIResult, err error) {
	var (
		localVarHTTPMethod   = nethttp.MethodGet
		localVarPostBody     interface{}
		localVarQueryParams  neturl.Values
		localVarHeaderParams = nethttp.Header{}
	)

	// Make sure Path Params are properly set
	if req.BillDate == nil || len(*req.BillDate) == 0 {
		return nil, nil, fmt.Errorf("field `BillDate` is required and must be specified in TradeBillRequest")
	}
	if req.BillType == nil || len(*req.BillType) == 0 {
		return nil, nil, fmt.Errorf("field `BillType` is required and must be specified in TradeBillRequest")
	}

	localVarPath := consts.WechatPayAPIServer + "/v3/bill/tradebill"

	// Setup Query Params
	localVarQueryParams = neturl.Values{}
	localVarQueryParams.Add("bill_date", core.ParameterToString(*req.BillDate, ""))
	localVarQueryParams.Add("bill_type", core.ParameterToString(*req.BillType, ""))

	// Determine the Content-Type Header
	var localVarHTTPContentTypes []string
	// Setup Content-Type
	localVarHTTPContentType := core.SelectHeaderContentType(localVarHTTPContentTypes)

	// Perform Http Request
	result, err = a.Client.Request(ctx, localVarHTTPMethod, localVarPath, localVarHeaderParams, localVarQueryParams, localVarPostBody, localVarHTTPContentType)
	if err != nil {
		return nil, result, err
	}

	// Extract bill.TradeBillResponse from Http Response
	resp = new(TradeBillResponse)
	err = core.UnMarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}
	return resp, result, nil
}

// DownloadBill 下载账单
func (a *BillApiService) DownloadBill(ctx context.Context, url string) ([]byte, error) {
	var (
		localVarHTTPMethod   = nethttp.MethodGet
		localVarPostBody     interface{}
		localVarQueryParams  neturl.Values
		localVarHeaderParams = nethttp.Header{}
	)

	localVarPath := url

	// Setup Query Params
	localVarQueryParams = neturl.Values{}

	// Determine the Content-Type Header
	localVarHTTPContentTypes := []string{}
	// Setup Content-Type
	localVarHTTPContentType := core.SelectHeaderContentType(localVarHTTPContentTypes)

	// 初始化不做校验的客户端
	newClient := core.NewClientWithValidator(a.Client, &validators.NullValidator{})
	// Perform Http Request
	result, err := newClient.Request(ctx, localVarHTTPMethod, localVarPath, localVarHeaderParams, localVarQueryParams, localVarPostBody, localVarHTTPContentType)
	if err != nil {
		return nil, err
	}
	httpResp := result.Response

	var body []byte
	body, err = io.ReadAll(httpResp.Body)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Errorf("failed to close body: %v", err)
		}
	}(httpResp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
