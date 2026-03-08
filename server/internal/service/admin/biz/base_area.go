package biz

import (
	"context"
	"strconv"
	"sync"

	"github.com/liujitcn/go-sdk/cache"
	"github.com/liujitcn/go-utils/str"
	"github.com/liujitcn/shop-admin/server/internal/data"
)

var codeMap map[string]string
var lock sync.RWMutex

type BaseAreaCase struct {
	data.BaseAreaRepo
	cache cache.Cache
}

// NewBaseAreaCase new a BaseArea use case.
func NewBaseAreaCase(
	baseAreaRepo data.BaseAreaRepo,
) *BaseAreaCase {
	return &BaseAreaCase{
		BaseAreaRepo: baseAreaRepo,
	}
}

func (c *BaseAreaCase) GetAddressListByCode(ctx context.Context, code string) []string {
	lock.RLock()
	defer lock.RUnlock()
	res := make([]string, 0)
	if codeMap == nil {
		list, err := c.FindAll(ctx, &data.BaseAreaCondition{})
		if err != nil {
			return res
		}
		codeMap = make(map[string]string)
		for _, item := range list {
			codeMap[strconv.FormatInt(item.ID, 10)] = item.Name
		}
	}
	codeList := str.ConvertJsonStringToStringArray(code)
	for _, item := range codeList {
		if v, ok := codeMap[item]; ok {
			res = append(res, v)
		} else {
			res = append(res, item)
		}
	}
	return res
}
