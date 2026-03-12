package core

import (
	"sync"
	"time"

	authData "github.com/liujitcn/kratos-kit/auth/data"
	"github.com/liujitcn/kratos-kit/bootstrap"
	"github.com/liujitcn/kratos-kit/cache"
	"github.com/liujitcn/kratos-kit/database/gorm"
	"github.com/liujitcn/kratos-kit/pprof"
	"github.com/liujitcn/kratos-kit/queue"
	"github.com/liujitcn/kratos-kit/sdk"
	"github.com/liujitcn/shop-admin/server/internal/service/file/biz"
)

type ShopCore struct {
	Ctx       *bootstrap.Context
	Cache     cache.Cache
	Queue     queue.Queue
	Gorm      *gorm.Client
	FileCase  *biz.FileCase
	UserToken *authData.UserToken
	quitChan  chan struct{} //退出Chan
	closeOnce sync.Once
	taskTimer *time.Timer
	rwLock    sync.RWMutex //异步数据锁
}

// NewShopCore create a device service core data struct.
func NewShopCore(
	ctx *bootstrap.Context,
	cache cache.Cache,
	queue queue.Queue,
	gorm *gorm.Client,
	pprof pprof.Pprof,
	userToken *authData.UserToken,
) (*ShopCore, func(), error) {

	// 设置全局变量
	sdk.Runtime.SetGormClient(gorm)
	sdk.Runtime.SetCache(cache)
	sdk.Runtime.SetQueue(queue)

	// 启动服务监控
	if pprof != nil {
		pprof.Start()
	}

	usc := ShopCore{
		Ctx:       ctx,
		Gorm:      gorm,
		Cache:     cache,
		Queue:     queue,
		UserToken: userToken,

		quitChan:  make(chan struct{}),
		closeOnce: sync.Once{},
		taskTimer: nil,
		rwLock:    sync.RWMutex{},
	}

	// 启动后台服务
	go usc.Serve()

	cleanup := func() {
		usc.Close()
	}
	return &usc, cleanup, nil
}

func (sc *ShopCore) Close() {
	sc.closeOnce.Do(func() {
		if sc.taskTimer != nil {
			sc.taskTimer.Stop()
		}
		close(sc.quitChan)
	})
}

// Serve 缓存加载和刷新线程
func (sc *ShopCore) Serve() {
	// 启动队列
	sc.Queue.Run()
	//循环处理同步事件
	for {
		select {
		case <-sc.quitChan:
			return
		}
	}
}
