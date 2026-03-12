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

type BaseDictItemCase struct {
	data.BaseDictItemRepo
	baseDictRepo data.BaseDictRepo
}

// NewBaseDictItemCase new a BaseDictItem use case.
func NewBaseDictItemCase(baseDictRepo data.BaseDictRepo, baseDictItemRepo data.BaseDictItemRepo) *BaseDictItemCase {
	return &BaseDictItemCase{
		baseDictRepo:     baseDictRepo,
		BaseDictItemRepo: baseDictItemRepo,
	}
}
func (c *BaseDictItemCase) GetFromID(ctx context.Context, id int64) (*models.BaseDictItem, error) {
	return c.Find(ctx, &data.BaseDictItemCondition{
		Id: id,
	})
}

func (c *BaseDictItemCase) Page(ctx context.Context, req *admin.PageBaseDictItemRequest) (*admin.PageBaseDictItemResponse, error) {
	condition := &data.BaseDictItemCondition{
		DictId: req.GetDictId(),
		Status: int32(req.GetStatus()),
		Label:  req.GetLabel(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.BaseDictItem, 0)
	for _, item := range page {
		list = append(list, &admin.BaseDictItem{
			Id:        item.ID,
			DictId:    item.DictID,
			Value:     item.Value,
			Label:     item.Label,
			TagType:   item.TagType,
			Sort:      item.Sort,
			Status:    common.Status(item.Status),
			CreatedAt: _time.TimeToTimeString(item.CreatedAt),
			UpdatedAt: _time.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageBaseDictItemResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseDictItemCase) List(ctx context.Context, condition *data.BaseDictItemCondition) ([]*models.BaseDictItem, error) {
	return c.FindAll(ctx, condition)
}

func (c *BaseDictItemCase) ConvertToProto(item *models.BaseDictItem) *admin.BaseDictItemForm {
	res := &admin.BaseDictItemForm{
		Id:      item.ID,
		DictId:  item.DictID,
		Value:   item.Value,
		Label:   item.Label,
		TagType: item.TagType,
		Sort:    item.Sort,
		Status:  trans.Enum(common.Status(item.Status)),
	}
	return res
}

func (c *BaseDictItemCase) ConvertToModel(item *admin.BaseDictItemForm) *models.BaseDictItem {
	res := &models.BaseDictItem{
		ID:      item.GetId(),
		DictID:  item.GetDictId(),
		Value:   item.GetValue(),
		Label:   item.GetLabel(),
		TagType: item.GetTagType(),
		Sort:    item.GetSort(),
		Status:  int32(item.GetStatus()),
	}
	return res
}
