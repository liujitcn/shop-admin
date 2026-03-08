package biz

import (
	"context"

	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type ShopServiceCase struct {
	data.ShopServiceRepo
}

// NewShopServiceCase new a ShopService use case.
func NewShopServiceCase(shopServiceRepo data.ShopServiceRepo) *ShopServiceCase {
	return &ShopServiceCase{
		ShopServiceRepo: shopServiceRepo,
	}
}

func (c *ShopServiceCase) Page(ctx context.Context, req *admin.PageShopServiceRequest) (*admin.PageShopServiceResponse, error) {
	condition := &data.ShopServiceCondition{
		Label:  req.GetLabel(),
		Status: int32(req.GetStatus()),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.ShopService, 0)
	for _, item := range page {
		list = append(list, &admin.ShopService{
			Id:        item.ID,
			Label:     item.Label,
			Value:     item.Value,
			Sort:      item.Sort,
			Status:    common.Status(item.Status),
			CreatedAt: timeutil.TimeToTimeString(item.CreatedAt),
			UpdatedAt: timeutil.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageShopServiceResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *ShopServiceCase) List(ctx context.Context, condition *data.ShopServiceCondition) ([]*models.ShopService, error) {
	return c.FindAll(ctx, condition)
}

func (c *ShopServiceCase) ConvertToProto(item *models.ShopService) *admin.ShopServiceForm {
	res := &admin.ShopServiceForm{
		Id:     item.ID,
		Label:  item.Label,
		Value:  item.Value,
		Sort:   item.Sort,
		Status: trans.Enum(common.Status(item.Status)),
	}
	return res
}

func (c *ShopServiceCase) ConvertToModel(item *admin.ShopServiceForm) *models.ShopService {
	res := &models.ShopService{
		ID:     item.GetId(),
		Label:  item.GetLabel(),
		Value:  item.GetValue(),
		Sort:   item.GetSort(),
		Status: int32(item.GetStatus()),
	}
	return res
}
