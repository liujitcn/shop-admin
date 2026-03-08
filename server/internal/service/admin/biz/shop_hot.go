package biz

import (
	"context"

	"github.com/liujitcn/go-utils/str"
	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type ShopHotCase struct {
	data.ShopHotRepo
}

// NewShopHotCase new a ShopHot use case.
func NewShopHotCase(shopHotRepo data.ShopHotRepo) *ShopHotCase {
	return &ShopHotCase{
		ShopHotRepo: shopHotRepo,
	}
}

func (c *ShopHotCase) GetFromID(ctx context.Context, id int64) (*models.ShopHot, error) {
	return c.Find(ctx, &data.ShopHotCondition{
		Id: id,
	})
}

func (c *ShopHotCase) List(ctx context.Context, condition *data.ShopHotCondition) ([]*models.ShopHot, error) {
	return c.FindAll(ctx, condition)
}

func (c *ShopHotCase) Page(ctx context.Context, req *admin.PageShopHotRequest) (*admin.PageShopHotResponse, error) {
	condition := &data.ShopHotCondition{
		Status: int32(req.GetStatus()),
		Title:  req.GetTitle(),
		Desc:   req.GetDesc(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.ShopHot, 0)
	for _, item := range page {
		list = append(list, &admin.ShopHot{
			Id:        item.ID,
			Title:     item.Title,
			Desc:      item.Desc,
			Sort:      item.Sort,
			Status:    common.Status(item.Status),
			CreatedAt: timeutil.TimeToTimeString(item.CreatedAt),
			UpdatedAt: timeutil.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageShopHotResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *ShopHotCase) ConvertToProto(item *models.ShopHot) *admin.ShopHotForm {
	return &admin.ShopHotForm{
		Id:      item.ID,
		Title:   item.Title,
		Desc:    item.Desc,
		Banner:  item.Banner,
		Picture: str.ConvertJsonStringToStringArray(item.Picture),
		Sort:    item.Sort,
		Status:  trans.Enum(common.Status(item.Status)),
	}
}

func (c *ShopHotCase) ConvertToModel(item *admin.ShopHotForm) *models.ShopHot {
	res := &models.ShopHot{
		ID:      item.GetId(),
		Title:   item.GetTitle(),
		Desc:    item.GetDesc(),
		Banner:  item.GetBanner(),
		Picture: str.ConvertStringArrayToString(item.GetPicture()),
		Sort:    item.GetSort(),
		Status:  int32(item.GetStatus()),
	}
	return res
}
