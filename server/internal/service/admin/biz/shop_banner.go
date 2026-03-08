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

type ShopBannerCase struct {
	data.ShopBannerRepo
}

// NewShopBannerCase new a ShopBanner use case.
func NewShopBannerCase(shopBannerRepo data.ShopBannerRepo) *ShopBannerCase {
	return &ShopBannerCase{
		ShopBannerRepo: shopBannerRepo,
	}
}

func (c *ShopBannerCase) Page(ctx context.Context, req *admin.PageShopBannerRequest) (*admin.PageShopBannerResponse, error) {
	condition := &data.ShopBannerCondition{
		Site:   int32(req.GetSite()),
		Type:   int32(req.GetType()),
		Status: int32(req.GetStatus()),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.ShopBanner, 0)
	for _, item := range page {
		list = append(list, &admin.ShopBanner{
			Id:        item.ID,
			Site:      common.ShopBannerSite(item.Site),
			Picture:   item.Picture,
			Type:      common.ShopBannerType(item.Type),
			Href:      item.Href,
			Sort:      item.Sort,
			Status:    common.Status(item.Status),
			CreatedAt: timeutil.TimeToTimeString(item.CreatedAt),
			UpdatedAt: timeutil.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageShopBannerResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *ShopBannerCase) List(ctx context.Context, condition *data.ShopBannerCondition) ([]*models.ShopBanner, error) {
	return c.FindAll(ctx, condition)
}

func (c *ShopBannerCase) ConvertToProto(item *models.ShopBanner) *admin.ShopBannerForm {
	res := &admin.ShopBannerForm{
		Id:      item.ID,
		Site:    trans.Enum(common.ShopBannerSite(item.Site)),
		Picture: item.Picture,
		Type:    trans.Enum(common.ShopBannerType(item.Type)),
		Href:    item.Href,
		Sort:    item.Sort,
		Status:  trans.Enum(common.Status(item.Status)),
	}
	return res
}

func (c *ShopBannerCase) ConvertToModel(item *admin.ShopBannerForm) *models.ShopBanner {
	res := &models.ShopBanner{
		ID:      item.GetId(),
		Site:    int32(item.GetSite()),
		Picture: item.GetPicture(),
		Type:    int32(item.GetType()),
		Href:    item.GetHref(),
		Sort:    item.GetSort(),
		Status:  int32(item.GetStatus()),
	}
	return res
}
