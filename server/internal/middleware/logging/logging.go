package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/status"
	"github.com/liujitcn/go-sdk"
	queueData "github.com/liujitcn/go-sdk/queue/data"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/file"
	"github.com/liujitcn/shop-admin/server/api/gen/go/login"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	"github.com/liujitcn/shop-admin/server/internal/service/admin/biz"
	"github.com/liujitcn/shop-gorm-gen/models"
	"github.com/mileusna/useragent"
	authnEngine "github.com/tx7do/kratos-authn/engine"
	"google.golang.org/grpc/codes"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server is an server logging middleware.
func Server(logger log.Logger,
	userCase *biz.BaseUserCase,
	authenticator authnEngine.Authenticator,
) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var kind string

			startTime := time.Now()
			// 日志信息
			baseLog := models.BaseLog{
				RequestTime: startTime,
				// default code
				StatusCode: int32(status.FromGRPCCode(codes.OK)),
			}
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				baseLog.Operation = info.Operation()
				var fullErr error
				if htr, htrOk := info.(*http.Transport); htrOk {
					baseLog.RequestID = getRequestId(htr.Request())
					// 文件上传不存请求内容
					if !(htr.Operation() == file.OperationFileServiceMultiUploadFile || htr.Operation() == file.OperationFileServiceUploadFile || htr.Operation() == file.OperationFileServiceDownloadFile) {
						baseLog.RequestBody = extractArgs(req)
					}

					headers := htr.RequestHeader()
					headersMap := make(map[string]string)
					for _, key := range headers.Keys() {
						headersMap[key] = htr.RequestHeader().Get(key)
					}
					var headersBytes []byte
					headersBytes, fullErr = json.Marshal(headersMap)
					if fullErr == nil {
						baseLog.RequestHeader = string(headersBytes)
					}

					clientIp := getClientRealIP(htr.Request())
					referer, _ := url.QueryUnescape(htr.RequestHeader().Get(HeaderKeyReferer))
					requestUri, _ := url.QueryUnescape(htr.Request().RequestURI)

					baseLog.Method = htr.Request().Method
					baseLog.Path = htr.PathTemplate()
					baseLog.Referer = trans.Ptr(referer)
					baseLog.RequestURI = trans.Ptr(requestUri)
					baseLog.Location = trans.Ptr(clientIpToLocation(clientIp))

					if htr.Operation() == admin.OperationAuthServiceLogin {
						var loginRequest login.LoginRequest
						if fullErr = json.Unmarshal([]byte(baseLog.RequestBody), &loginRequest); fullErr == nil {
							userName := loginRequest.GetUserName()
							if len(userName) > 0 {
								baseLog.UserName = userName
								var baseUser *models.BaseUser
								baseUser, fullErr = userCase.GetFromUserName(htr.Request().Context(), userName)
								if fullErr == nil {
									baseLog.UserID = baseUser.ID
								}
							}
						}
					} else {
						authToken := htr.RequestHeader().Get(HeaderKeyAuthorization)
						ut := extractAuthToken(authToken, authenticator)
						if ut != nil {
							baseLog.UserID = ut.UserId
							baseLog.UserName = ut.UserName
						}
					}

					// 用户代理信息
					strUserAgent := htr.RequestHeader().Get(HeaderKeyUserAgent)
					ua := useragent.Parse(strUserAgent)

					var deviceName string
					if ua.Device != "" {
						deviceName = ua.Device
					} else {
						if ua.Desktop {
							deviceName = "PC"
						}
					}

					baseLog.UserAgent = ua.String
					baseLog.BrowserVersion = ua.Version
					baseLog.BrowserName = ua.Name
					baseLog.OsName = ua.OS
					baseLog.OsVersion = ua.OSVersion
					baseLog.ClientName = deviceName
				}
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				baseLog.StatusCode = se.Code
				baseLog.Reason = se.Reason
			}
			baseLog.CostTime = time.Since(startTime).Milliseconds()
			responseBytes, responseErr := json.Marshal(reply)
			if responseErr == nil {
				baseLog.Response = string(responseBytes)
			}
			level, stack := extractError(err)
			if len(stack) > 0 {
				baseLog.Reason = fmt.Sprintf("[%s]%s", baseLog.Reason, stack)
			}
			// 写入日志
			writeOperationLog(&baseLog)
			log.NewHelper(log.WithContext(ctx, logger)).Log(level,
				"kind", "server",
				"component", kind,
				"operation", baseLog.Operation,
				"args", baseLog.RequestBody,
				"code", baseLog.StatusCode,
				"reason", baseLog.Reason,
				"stack", stack,
				"latency", baseLog.CostTime,
			)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}

// 写入操作日志
func writeOperationLog(
	baseLog *models.BaseLog,
) {
	var err error
	// 加入日志队列
	q := sdk.Runtime.GetQueue()
	if q != nil {
		m := make(map[string]interface{})
		m["data"] = baseLog
		var message queueData.Message
		message, err = sdk.Runtime.GetStreamMessage(_const.Log, m)
		if err != nil {
			log.Errorf("GetStreamMessage error, %s", err.Error())
			//日志报错错误，不中断请求
		} else {
			err = q.Append(_const.Log, message)
			if err != nil {
				log.Errorf("Append message error, %s", err.Error())
			}
		}
	}
}
