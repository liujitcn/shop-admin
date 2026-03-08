package job

import (
	"strings"
	"time"

	"github.com/liujitcn/go-utils/str"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	"github.com/liujitcn/shop-admin/server/internal/sdk"
	"github.com/liujitcn/shop-admin/server/internal/service/admin/task"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type ExecJob struct {
	JobId        int64             // 任务ID
	Args         map[string]string // 任务参数
	InvokeTarget task.TaskExec
	Status       common.BaseJobLogStatus
	ErrMsg       string
}

// Run 函数任务执行
func (e *ExecJob) Run() {
	// 记录日志
	baseJobLog := models.BaseJobLog{
		JobID:       e.JobId,                            // 定时任务id
		Input:       str.ConvertAnyToJsonString(e.Args), // 任务参数
		ExecuteTime: time.Now(),                         // 执行时间
	}
	ret, err := e.InvokeTarget.Exec(e.Args)
	if err != nil {
		e.Status = common.BaseJobLogStatus_FAIL
		e.ErrMsg = err.Error()
	} else {
		e.Status = common.BaseJobLogStatus_SUCCESS
	}
	// 执行结果
	baseJobLog.Output = strings.Join(ret, "<br/>")
	// 执行结果-成功
	baseJobLog.Status = int32(e.Status)
	baseJobLog.Error = e.ErrMsg
	// 执行时间
	baseJobLog.ProcessTime = int32(time.Now().Sub(baseJobLog.ExecuteTime).Milliseconds())
	// 加入日志队列
	sdk.AddQueue(_const.JobLog, baseJobLog)
	return
}
