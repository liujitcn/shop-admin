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

type BaseDictCase struct {
	data.BaseDictRepo
}

// NewBaseDictCase new a BaseDict use case.
func NewBaseDictCase(baseDictRepo data.BaseDictRepo) *BaseDictCase {
	return &BaseDictCase{
		BaseDictRepo: baseDictRepo,
	}
}

func (c *BaseDictCase) GetFromID(ctx context.Context, id int64) (*models.BaseDict, error) {
	return c.Find(ctx, &data.BaseDictCondition{
		Id: id,
	})
}

func (c *BaseDictCase) List(ctx context.Context, condition *data.BaseDictCondition) ([]*models.BaseDict, error) {
	return c.FindAll(ctx, condition)
}

func (c *BaseDictCase) Page(ctx context.Context, req *admin.PageBaseDictRequest) (*admin.PageBaseDictResponse, error) {
	condition := &data.BaseDictCondition{
		Status: int32(req.GetStatus()),
		Name:   req.GetName(),
		Code:   req.GetCode(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.BaseDict, 0)
	for _, item := range page {
		list = append(list, &admin.BaseDict{
			Id:        item.ID,
			Code:      item.Code,
			Name:      item.Name,
			Status:    common.Status(item.Status),
			CreatedAt: _time.TimeToTimeString(item.CreatedAt),
			UpdatedAt: _time.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageBaseDictResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseDictCase) ConvertToProto(item *models.BaseDict) *admin.BaseDictForm {
	return &admin.BaseDictForm{
		Id:     item.ID,
		Name:   item.Name,
		Code:   item.Code,
		Status: trans.Enum(common.Status(item.Status)),
	}
}

func (c *BaseDictCase) ConvertToModel(item *admin.BaseDictForm) *models.BaseDict {
	res := &models.BaseDict{
		ID:     item.GetId(),
		Code:   item.GetCode(),
		Name:   item.GetName(),
		Status: int32(item.GetStatus()),
	}
	return res
}
