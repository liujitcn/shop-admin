package biz

import (
	"context"

	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type ShopHotItemCase struct {
	tx genData.Transaction
	data.ShopHotItemRepo
	shopHotRepo      data.ShopHotRepo
	shopHotGoodsRepo data.ShopHotGoodsRepo
}

// NewShopHotItemCase new a ShopHotItem use case.
func NewShopHotItemCase(
	tx genData.Transaction,
	shopHotRepo data.ShopHotRepo,
	shopHotItemRepo data.ShopHotItemRepo,
	shopHotGoodsRepo data.ShopHotGoodsRepo,
) *ShopHotItemCase {
	return &ShopHotItemCase{
		tx:               tx,
		shopHotRepo:      shopHotRepo,
		ShopHotItemRepo:  shopHotItemRepo,
		shopHotGoodsRepo: shopHotGoodsRepo,
	}
}
func (c *ShopHotItemCase) GetFromID(ctx context.Context, id int64) (*admin.ShopHotItemForm, error) {
	shopHotItem, err := c.Find(ctx, &data.ShopHotItemCondition{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	res := c.ConvertToProto(shopHotItem)
	res.GoodsIds, err = c.shopHotGoodsRepo.FindAll(ctx, shopHotItem.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ShopHotItemCase) Page(ctx context.Context, req *admin.PageShopHotItemRequest) (*admin.PageShopHotItemResponse, error) {
	condition := &data.ShopHotItemCondition{
		HotId:  req.GetHotId(),
		Status: int32(req.GetStatus()),
		Title:  req.GetTitle(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.ShopHotItem, 0)
	for _, item := range page {
		list = append(list, &admin.ShopHotItem{
			Id:        item.ID,
			HotId:     item.HotID,
			Title:     item.Title,
			Sort:      item.Sort,
			Status:    common.Status(item.Status),
			CreatedAt: timeutil.TimeToTimeString(item.CreatedAt),
			UpdatedAt: timeutil.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageShopHotItemResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *ShopHotItemCase) List(ctx context.Context, condition *data.ShopHotItemCondition) ([]*models.ShopHotItem, error) {
	return c.FindAll(ctx, condition)
}

func (c *ShopHotItemCase) Create1(ctx context.Context, req *admin.ShopHotItemForm) error {
	shopHotItem := c.ConvertToModel(req)
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err := c.Create(ctx, shopHotItem)
		if err != nil {
			return err
		}
		return c.shopHotGoodsRepo.Create(ctx, shopHotItem.ID, req.GetGoodsIds())
	})
}

func (c *ShopHotItemCase) Update(ctx context.Context, req *admin.ShopHotItemForm) error {
	shopHotItem := c.ConvertToModel(req)
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err := c.UpdateByID(ctx, shopHotItem)
		if err != nil {
			return err
		}
		err = c.shopHotGoodsRepo.Delete(ctx, shopHotItem.ID)
		if err != nil {
			return err
		}
		return c.shopHotGoodsRepo.Create(ctx, shopHotItem.ID, req.GetGoodsIds())
	})
}

func (c *ShopHotItemCase) ConvertToProto(item *models.ShopHotItem) *admin.ShopHotItemForm {
	res := &admin.ShopHotItemForm{
		Id:     item.ID,
		HotId:  item.HotID,
		Title:  item.Title,
		Sort:   item.Sort,
		Status: trans.Enum(common.Status(item.Status)),
	}
	return res
}

func (c *ShopHotItemCase) ConvertToModel(item *admin.ShopHotItemForm) *models.ShopHotItem {
	res := &models.ShopHotItem{
		ID:     item.GetId(),
		HotID:  item.GetHotId(),
		Title:  item.GetTitle(),
		Sort:   item.GetSort(),
		Status: int32(item.GetStatus()),
	}
	return res
}
