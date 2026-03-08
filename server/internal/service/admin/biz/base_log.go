package biz

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	queueData "github.com/liujitcn/go-sdk/queue/data"
	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type BaseLogCase struct {
	data.BaseLogRepo
}

// NewBaseLogCase new a BaseLog use case.
func NewBaseLogCase(baseLogRepo data.BaseLogRepo) *BaseLogCase {
	return &BaseLogCase{
		BaseLogRepo: baseLogRepo,
	}
}

func (c *BaseLogCase) GetFromID(ctx context.Context, id int64) (*models.BaseLog, error) {
	return c.Find(ctx, &data.BaseLogCondition{
		Id: id,
	})
}

func (c *BaseLogCase) Page(ctx context.Context, req *admin.PageBaseLogRequest) (*admin.PageBaseLogResponse, error) {
	requestTime := req.GetRequestTime()
	var startTime, endTime *time.Time
	if len(requestTime) == 2 {
		startTime = timeutil.StringTimeToTime(requestTime[0])
		endTime = timeutil.StringTimeToTime(requestTime[1])
		if endTime != nil {
			t := endTime.AddDate(0, 0, 1)
			endTime = &t
		}
	}
	condition := &data.BaseLogCondition{
		Operation:        req.GetOperation(),
		StatusCode:       req.GetStatusCode(),
		RequestStartTime: startTime,
		RequestEndTime:   endTime,
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.BaseLog, 0)
	for _, item := range page {
		list = append(list, c.ConvertToProto(item))
	}

	return &admin.PageBaseLogResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseLogCase) ConvertToProto(item *models.BaseLog) *admin.BaseLog {
	costTime := time.Duration(item.CostTime) * time.Millisecond
	return &admin.BaseLog{
		Id:             item.ID,
		RequestId:      item.RequestID,
		RequestTime:    timeutil.TimeToTimeString(item.RequestTime),
		Method:         item.Method,
		Operation:      item.Operation,
		Path:           item.Path,
		Referer:        item.Referer,
		RequestUri:     item.RequestURI,
		RequestHeader:  item.RequestHeader,
		RequestBody:    item.RequestBody,
		Response:       item.Response,
		CostTime:       costTime.String(),
		Success:        item.Success,
		StatusCode:     item.StatusCode,
		Reason:         item.Reason,
		Location:       item.Location,
		UserId:         item.UserID,
		UserName:       item.UserName,
		ClientIp:       item.ClientIP,
		UserAgent:      item.UserAgent,
		BrowserName:    item.BrowserName,
		BrowserVersion: item.BrowserVersion,
		ClientId:       item.ClientID,
		ClientName:     item.ClientName,
		OsName:         item.OsName,
		OsVersion:      item.OsVersion,
	}
}

func (c *BaseLogCase) SaveLog(message queueData.Message) error {
	rb, err := json.Marshal(message.Values)
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		return err
	}
	var m map[string]*models.BaseLog
	err = json.Unmarshal(rb, &m)
	if err != nil {
		log.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	if v, ok := m["data"]; ok {
		err = c.Create(context.TODO(), v)
		if err != nil {
			return err
		}
	}
	return nil
}
