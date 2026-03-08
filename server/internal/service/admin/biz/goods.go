package biz

import (
	"context"
	"fmt"

	"github.com/liujitcn/go-utils/str"
	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type GoodsCase struct {
	tx genData.Transaction
	data.GoodsRepo
	goodsCategoryRepo data.GoodsCategoryRepo
	goodsPropCase     *GoodsPropCase
	goodsSpecCase     *GoodsSpecCase
	goodsSkuCase      *GoodsSkuCase
}

// NewGoodsCase new a Goods use case.
func NewGoodsCase(
	tx genData.Transaction,
	goodsRepo data.GoodsRepo,
	goodsCategoryRepo data.GoodsCategoryRepo,
	goodsPropCase *GoodsPropCase,
	goodsSpecCase *GoodsSpecCase,
	goodsSkuCase *GoodsSkuCase,
) *GoodsCase {
	return &GoodsCase{
		tx:                tx,
		GoodsRepo:         goodsRepo,
		goodsCategoryRepo: goodsCategoryRepo,
		goodsPropCase:     goodsPropCase,
		goodsSpecCase:     goodsSpecCase,
		goodsSkuCase:      goodsSkuCase,
	}
}
func (c *GoodsCase) GetFromID(ctx context.Context, id int64) (*admin.GoodsForm, error) {
	goods, err := c.Find(ctx, &data.GoodsCondition{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	goodsForm := c.ConvertToProto(goods)

	// 查询分类路径
	var category *models.GoodsCategory
	category, err = c.goodsCategoryRepo.Find(ctx, &data.GoodsCategoryCondition{Id: goodsForm.GetCategoryId()})
	if err != nil {
		return nil, err
	}
	var parentCategory *models.GoodsCategory
	parentCategory, err = c.goodsCategoryRepo.Find(ctx, &data.GoodsCategoryCondition{Id: category.ParentID})
	if err != nil {
		return nil, err
	}
	goodsForm.CategoryName = fmt.Sprintf("%s/%s", parentCategory.Name, category.Name)

	// 属性
	goodsForm.PropList, err = c.goodsPropCase.ListByGoodsId(ctx, goodsForm.Id)
	if err != nil {
		return nil, err
	}
	// 规格
	goodsForm.SpecList, err = c.goodsSpecCase.ListByGoodsId(ctx, goodsForm.Id)
	if err != nil {
		return nil, err
	}
	// SKU
	goodsForm.SkuList, err = c.goodsSkuCase.ListByGoodsId(ctx, goodsForm.Id)
	if err != nil {
		return nil, err
	}
	return goodsForm, nil
}

func (c *GoodsCase) Page(ctx context.Context, req *admin.PageGoodsRequest) (*admin.PageGoodsResponse, error) {
	condition := &data.GoodsCondition{
		Name:       req.GetName(),
		CategoryId: req.GetCategoryId(),
		Status:     int32(req.GetStatus()),
	}
	if condition.CategoryId > 0 {
		// 查询分类路径
		category, err := c.goodsCategoryRepo.Find(ctx, &data.GoodsCategoryCondition{Id: req.GetCategoryId()})
		if err != nil {
			return nil, err
		}
		if category.ParentID == 0 {
			condition.CategoryId = 0
			condition.CategoryPath = fmt.Sprintf("%s/", category.Path)
		}
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.Goods, 0)
	for _, item := range page {
		list = append(list, &admin.Goods{
			Id:            item.ID,
			CategoryId:    item.CategoryID,
			Name:          item.Name,
			Desc:          item.Desc,
			Picture:       item.Picture,
			Price:         item.Price,
			DiscountPrice: item.DiscountPrice,
			InitSaleNum:   item.InitSaleNum,
			RealSaleNum:   item.RealSaleNum,
			Status:        common.GoodsStatus(item.Status),
			CreatedAt:     timeutil.TimeToTimeString(item.CreatedAt),
			UpdatedAt:     timeutil.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageGoodsResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *GoodsCase) Delete(ctx context.Context, ids []int64) error {
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err := c.GoodsRepo.Delete(ctx, ids)
		if err != nil {
			return err
		}
		for _, item := range ids {
			// 删除
			err = c.goodsPropCase.DeleteByGoodsId(ctx, item)
			if err != nil {
				return err
			}
			err = c.goodsSpecCase.DeleteByGoodsId(ctx, item)
			if err != nil {
				return err
			}
			err = c.goodsSkuCase.DeleteByGoodsId(ctx, item)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *GoodsCase) Create(ctx context.Context, goodsForm *admin.GoodsForm) error {
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		// 基础信息
		goods := c.ConvertToModel(goodsForm)
		skuList := goodsForm.GetSkuList()
		for idx, sku := range skuList {
			if idx == 0 {
				goods.Price = sku.Price
				goods.DiscountPrice = sku.DiscountPrice
			}
			goods.InitSaleNum += sku.InitSaleNum
			goods.RealSaleNum += sku.RealSaleNum
		}
		err := c.GoodsRepo.Create(ctx, goods)
		if err != nil {
			return err
		}
		// 属性
		err = c.goodsPropCase.BatchCreate(ctx, goods.ID, goodsForm.GetPropList())
		if err != nil {
			return err
		}
		// 规格
		err = c.goodsSpecCase.BatchCreate(ctx, goods.ID, goodsForm.GetSpecList())
		if err != nil {
			return err
		}
		// SKU
		err = c.goodsSkuCase.BatchCreate(ctx, goods.ID, skuList)
		if err != nil {
			return err
		}

		return nil
	})
}

func (c *GoodsCase) Update(ctx context.Context, goodsForm *admin.GoodsForm) error {
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		// 基础信息
		goods := c.ConvertToModel(goodsForm)
		skuList := goodsForm.GetSkuList()
		for idx, sku := range skuList {
			if idx == 0 {
				goods.Price = sku.Price
				goods.DiscountPrice = sku.DiscountPrice
			}
			goods.InitSaleNum += sku.InitSaleNum
			goods.RealSaleNum += sku.RealSaleNum
		}
		err := c.GoodsRepo.UpdateByID(ctx, goods)
		if err != nil {
			return err
		}
		// 属性
		err = c.goodsPropCase.BatchCreate(ctx, goods.ID, goodsForm.GetPropList())
		if err != nil {
			return err
		}
		// 规格
		err = c.goodsSpecCase.BatchCreate(ctx, goods.ID, goodsForm.GetSpecList())
		if err != nil {
			return err
		}
		// SKU
		err = c.goodsSkuCase.BatchCreate(ctx, goods.ID, skuList)
		if err != nil {
			return err
		}

		return nil
	})
}

func (c *GoodsCase) List(ctx context.Context, condition *data.GoodsCondition) ([]*models.Goods, error) {
	return c.FindAll(ctx, condition)
}

func (c *GoodsCase) ConvertToProto(item *models.Goods) *admin.GoodsForm {
	res := &admin.GoodsForm{
		Id:         item.ID,
		CategoryId: trans.Int64(item.CategoryID),
		Name:       item.Name,
		Desc:       item.Desc,
		Picture:    item.Picture,
		Banner:     str.ConvertJsonStringToStringArray(item.Banner),
		Detail:     str.ConvertJsonStringToStringArray(item.Detail),
		Status:     trans.Enum(common.GoodsStatus(item.Status)),
	}
	return res
}

func (c *GoodsCase) ConvertToModel(item *admin.GoodsForm) *models.Goods {
	res := &models.Goods{
		ID:            item.GetId(),
		CategoryID:    item.GetCategoryId(),
		Name:          item.GetName(),
		Desc:          item.GetDesc(),
		Picture:       item.GetPicture(),
		Banner:        str.ConvertStringArrayToString(item.GetBanner()),
		Detail:        str.ConvertStringArrayToString(item.GetDetail()),
		Price:         0,
		DiscountPrice: 0,
		InitSaleNum:   0,
		RealSaleNum:   0,
		Status:        int32(item.GetStatus()),
	}
	return res
}
