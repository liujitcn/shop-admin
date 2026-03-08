package biz

import (
	"context"
	"errors"

	"github.com/liujitcn/shop-admin/server/api/gen/go/config"
	"github.com/liujitcn/shop-admin/server/internal/data"
)

type ConfigCase struct {
	data.BaseConfigRepo
}

// NewConfigCase new a Config use case.
func NewConfigCase(baseConfigRepo data.BaseConfigRepo) *ConfigCase {
	return &ConfigCase{
		BaseConfigRepo: baseConfigRepo,
	}
}

func (c *ConfigCase) GetConfig(ctx context.Context, req *config.ConfigRequest) (*config.ConfigResponse, error) {
	if req.GetSite() == 0 {
		return nil, errors.New("位置不能为空")
	}
	list, err := c.FindAll(ctx, &data.BaseConfigCondition{
		Site: int32(req.GetSite()),
	})
	if err != nil {
		return nil, err
	}
	resData := make([]*config.ConfigResponse_Data, 0)
	for _, item := range list {
		resData = append(resData, &config.ConfigResponse_Data{
			Key:   item.Key,
			Value: item.Value,
		})
	}
	return &config.ConfigResponse{
		Data: resData,
	}, nil
}
