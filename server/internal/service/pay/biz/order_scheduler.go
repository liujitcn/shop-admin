package biz

import (
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type OrderSchedulerCase struct {
	timers sync.Map // 存储订单ID和对应的定时器
}

func NewOrderSchedulerCase() *OrderSchedulerCase {
	return &OrderSchedulerCase{}
}

func (s *OrderSchedulerCase) DeleteScheduled(orderId int64) {
	if timer, ok := s.timers.Load(orderId); ok {
		timer.(*time.Timer).Stop()
		log.Infof("order schedule delete %d", orderId)
		s.timers.Delete(orderId)
	}
}
