package biz

import (
	"context"

	_string "github.com/liujitcn/go-utils/string"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type GoodsSkuCase struct {
	data.GoodsSkuRepo
}

// NewGoodsSkuCase new a GoodsSku use case.
func NewGoodsSkuCase(goodsSkuRepo data.GoodsSkuRepo) *GoodsSkuCase {
	return &GoodsSkuCase{
		GoodsSkuRepo: goodsSkuRepo,
	}
}
func (c *GoodsSkuCase) GetFromID(ctx context.Context, id int64) (*models.GoodsSku, error) {
	return c.Find(ctx, &data.GoodsSkuCondition{
		Id: id,
	})
}

func (c *GoodsSkuCase) Page(ctx context.Context, req *admin.PageGoodsSkuRequest) (*admin.PageGoodsSkuResponse, error) {
	condition := &data.GoodsSkuCondition{
		GoodsId: req.GetGoodsId(),
		SkuCode: req.GetSkuCode(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.GoodsSku, 0)
	for _, item := range page {
		list = append(list, c.ConvertToProto(item))
	}

	return &admin.PageGoodsSkuResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *GoodsSkuCase) BatchCreate(ctx context.Context, goodsId int64, sku []*admin.GoodsSku) error {
	if len(sku) == 0 {
		return nil
	}
	// 查询旧数据
	oldSkuList, err := c.FindAll(ctx, &data.GoodsSkuCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return err
	}
	// 规格编号和id map
	oldSkuIdMap := make(map[string]int64)
	for _, oldSku := range oldSkuList {
		oldSkuIdMap[oldSku.SkuCode] = oldSku.ID
	}

	skuList := make([]*models.GoodsSku, 0)
	for _, item := range sku {
		if id, ok := oldSkuIdMap[item.SkuCode]; ok {
			item.Id = id
			err = c.UpdateByID(ctx, c.ConvertToModel(item))
			if err != nil {
				return err
			}
			// 删除map
			delete(oldSkuIdMap, item.SkuCode)
		} else {
			item.GoodsId = goodsId
			skuList = append(skuList, c.ConvertToModel(item))
		}
	}
	if len(oldSkuIdMap) > 0 {
		oldSkuId := make([]int64, 0)
		for _, v := range oldSkuIdMap {
			oldSkuId = append(oldSkuId, v)
		}
		err = c.Delete(ctx, oldSkuId)
		if err != nil {
			return err
		}
	}
	if len(skuList) >= 0 {
		return c.GoodsSkuRepo.BatchCreate(ctx, skuList)
	}
	return nil
}

func (c *GoodsSkuCase) ListByGoodsId(ctx context.Context, goodsId int64) ([]*admin.GoodsSku, error) {
	all, err := c.FindAll(ctx, &data.GoodsSkuCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*admin.GoodsSku, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item))
	}
	return list, nil
}

func (c *GoodsSkuCase) ConvertToProto(item *models.GoodsSku) *admin.GoodsSku {
	res := &admin.GoodsSku{
		Id:            item.ID,
		GoodsId:       item.GoodsID,
		Picture:       item.Picture,
		SpecItem:      _string.ConvertJsonStringToStringArray(item.SpecItem),
		SkuCode:       item.SkuCode,
		Price:         item.Price,
		DiscountPrice: item.DiscountPrice,
		InitSaleNum:   item.InitSaleNum,
		RealSaleNum:   item.RealSaleNum,
		Inventory:     item.Inventory,
	}
	return res
}

func (c *GoodsSkuCase) ConvertToModel(item *admin.GoodsSku) *models.GoodsSku {
	res := &models.GoodsSku{
		ID:            item.GetId(),
		GoodsID:       item.GetGoodsId(),
		Picture:       item.GetPicture(),
		SpecItem:      _string.ConvertStringArrayToString(item.SpecItem),
		SkuCode:       item.GetSkuCode(),
		Price:         item.GetPrice(),
		DiscountPrice: item.GetDiscountPrice(),
		InitSaleNum:   item.GetInitSaleNum(),
		RealSaleNum:   item.GetRealSaleNum(),
		Inventory:     item.GetInventory(),
	}
	return res
}
