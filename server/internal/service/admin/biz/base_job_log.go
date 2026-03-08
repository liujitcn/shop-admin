package biz

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	queueData "github.com/liujitcn/go-sdk/queue/data"
	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type BaseJobLogCase struct {
	data.BaseJobLogRepo
}

// NewBaseJobLogCase new a BaseJobLog use case.
func NewBaseJobLogCase(baseJobLogRepo data.BaseJobLogRepo) *BaseJobLogCase {
	return &BaseJobLogCase{
		BaseJobLogRepo: baseJobLogRepo,
	}
}

func (c *BaseJobLogCase) GetFromID(ctx context.Context, id int64) (*models.BaseJobLog, error) {
	return c.Find(ctx, &data.BaseJobLogCondition{
		Id: id,
	})
}

func (c *BaseJobLogCase) Page(ctx context.Context, req *admin.PageBaseJobLogRequest) (*admin.PageBaseJobLogResponse, error) {
	executeTime := req.GetExecuteTime()
	var startTime, endTime *time.Time
	if len(executeTime) == 2 {
		startTime = timeutil.StringTimeToTime(executeTime[0])
		endTime = timeutil.StringTimeToTime(executeTime[1])
		if endTime != nil {
			t := endTime.AddDate(0, 0, 1)
			endTime = &t
		}
	}
	condition := &data.BaseJobLogCondition{
		JobId:            req.GetJobId(),
		Status:           int32(req.GetStatus()),
		ExecuteStartTime: startTime,
		ExecuteEndTime:   endTime,
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.BaseJobLog, 0)
	for _, item := range page {
		list = append(list, c.ConvertToProto(item))
	}

	return &admin.PageBaseJobLogResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseJobLogCase) ConvertToProto(item *models.BaseJobLog) *admin.BaseJobLog {
	processTime := time.Duration(item.ProcessTime) * time.Millisecond
	return &admin.BaseJobLog{
		Id:          item.ID,
		JobId:       item.JobID,
		Input:       item.Input,
		Output:      item.Output,
		Error:       item.Error,
		Status:      common.BaseJobLogStatus(item.Status),
		ProcessTime: processTime.String(),
		ExecuteTime: timeutil.TimeToTimeString(item.ExecuteTime),
	}
}

func (c *BaseJobLogCase) SaveJobLog(message queueData.Message) error {
	rb, err := json.Marshal(message.Values)
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		return err
	}
	var m map[string]*models.BaseJobLog
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
