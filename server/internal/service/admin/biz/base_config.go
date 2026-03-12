package biz

import (
	"context"

	_time "github.com/liujitcn/go-utils/time"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type BaseConfigCase struct {
	data.BaseConfigRepo
}

// NewBaseConfigCase new a BaseConfig use case.
func NewBaseConfigCase(baseConfigRepo data.BaseConfigRepo) *BaseConfigCase {
	return &BaseConfigCase{
		BaseConfigRepo: baseConfigRepo,
	}
}
func (c *BaseConfigCase) GetFromID(ctx context.Context, id int64) (*models.BaseConfig, error) {
	return c.Find(ctx, &data.BaseConfigCondition{
		Id: id,
	})
}

func (c *BaseConfigCase) Page(ctx context.Context, req *admin.PageBaseConfigRequest) (*admin.PageBaseConfigResponse, error) {
	condition := &data.BaseConfigCondition{
		Site:   int32(req.GetSite()),
		Name:   req.GetName(),
		Type:   int32(req.GetType()),
		Key:    req.GetKey(),
		Status: int32(req.GetStatus()),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.BaseConfig, 0)
	for _, item := range page {
		list = append(list, &admin.BaseConfig{
			Id:        item.ID,
			Site:      common.BaseConfigSite(item.Site),
			Name:      item.Name,
			Type:      common.BaseConfigType(item.Type),
			Key:       item.Key,
			Value:     item.Value,
			Status:    common.Status(item.Status),
			CreatedAt: _time.TimeToTimeString(item.CreatedAt),
			UpdatedAt: _time.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageBaseConfigResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseConfigCase) List(ctx context.Context, condition *data.BaseConfigCondition) ([]*models.BaseConfig, error) {
	return c.FindAll(ctx, condition)
}

func (c *BaseConfigCase) ConvertToProto(item *models.BaseConfig) *admin.BaseConfigForm {
	res := &admin.BaseConfigForm{
		Id:     item.ID,
		Site:   trans.Enum(common.BaseConfigSite(item.Site)),
		Name:   item.Name,
		Type:   trans.Enum(common.BaseConfigType(item.Type)),
		Key:    item.Key,
		Value:  item.Value,
		Status: trans.Enum(common.Status(item.Status)),
	}
	return res
}

func (c *BaseConfigCase) ConvertToModel(item *admin.BaseConfigForm) *models.BaseConfig {
	res := &models.BaseConfig{
		ID:     item.GetId(),
		Site:   int32(item.GetSite()),
		Name:   item.GetName(),
		Type:   int32(item.GetType()),
		Key:    item.GetKey(),
		Value:  item.GetValue(),
		Status: int32(item.GetStatus()),
	}
	return res
}
