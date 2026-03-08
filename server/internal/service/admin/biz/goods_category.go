package biz

import (
	"context"
	"strconv"
	"strings"

	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type GoodsCategoryCase struct {
	data.GoodsCategoryRepo
}

// NewGoodsCategoryCase new a GoodsCategory use case.
func NewGoodsCategoryCase(goodsCategoryRepo data.GoodsCategoryRepo) *GoodsCategoryCase {
	return &GoodsCategoryCase{
		GoodsCategoryRepo: goodsCategoryRepo,
	}
}

func (c *GoodsCategoryCase) GetFromID(ctx context.Context, id int64) (*models.GoodsCategory, error) {
	return c.Find(ctx, &data.GoodsCategoryCondition{
		Id: id,
	})
}

func (c *GoodsCategoryCase) NameMap(ctx context.Context, condition *data.GoodsCategoryCondition) map[int64]string {
	categoryNameMap := make(map[int64]string)
	categoryPathMap := make(map[int64]string)
	categoryList, err := c.FindAll(ctx, condition)
	if err == nil {
		for _, category := range categoryList {
			categoryNameMap[category.ID] = category.Name
			categoryPathMap[category.ID] = category.Path
		}
	}
	for categoryId, path := range categoryPathMap {
		paths := strings.Split(path, "/")
		pathName := make([]string, 0)
		for _, item := range paths {
			id, _ := strconv.ParseInt(item, 10, 64)
			if name, ok := categoryNameMap[id]; ok {
				pathName = append(pathName, name)
			}
		}
		categoryPathMap[categoryId] = strings.Join(pathName, "/")
	}
	return categoryPathMap
}

func (c *GoodsCategoryCase) Tree(ctx context.Context, condition *data.GoodsCategoryCondition) (*admin.TreeGoodsCategoryResponse, error) {
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &admin.TreeGoodsCategoryResponse{
		List: c.buildTree(list, 0),
	}, nil
}

func (c *GoodsCategoryCase) Option(ctx context.Context, condition *data.GoodsCategoryCondition) (*common.TreeOptionResponse, error) {
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &common.TreeOptionResponse{
		List: c.buildOption(list, 0, condition.ParentId == nil),
	}, nil
}

// buildTree 构建部门树状
func (c *GoodsCategoryCase) buildTree(categoryList []*models.GoodsCategory, parentId int64) []*admin.GoodsCategory {
	var res []*admin.GoodsCategory
	for _, item := range categoryList {
		if item.ParentID == parentId {
			category := &admin.GoodsCategory{
				Id:        item.ID,
				ParentId:  item.ParentID,
				Name:      item.Name,
				Picture:   item.Picture,
				Sort:      item.Sort,
				Status:    common.Status(item.Status),
				CreatedAt: timeutil.TimeToTimeString(item.CreatedAt),
				UpdatedAt: timeutil.TimeToTimeString(item.UpdatedAt),
			}
			category.Children = c.buildTree(categoryList, item.ID)
			res = append(res, category)
		}
	}
	return res
}

// buildTree 构建部门树形选择
func (c *GoodsCategoryCase) buildOption(categoryList []*models.GoodsCategory, parentId int64, disabled bool) []*common.TreeOptionResponse_Option {
	var res []*common.TreeOptionResponse_Option
	for _, item := range categoryList {
		if item.ParentID == parentId {
			category := &common.TreeOptionResponse_Option{
				Label:    item.Name,
				Value:    item.ID,
				Disabled: disabled && item.ParentID == 0,
			}
			category.Children = c.buildOption(categoryList, item.ID, disabled)
			res = append(res, category)
		}
	}
	return res
}

func (c *GoodsCategoryCase) ConvertToProto(item *models.GoodsCategory) *admin.GoodsCategoryForm {
	res := &admin.GoodsCategoryForm{
		Id:       item.ID,
		ParentId: trans.Int64(item.ParentID),
		Name:     item.Name,
		Picture:  item.Picture,
		Sort:     item.Sort,
		Status:   trans.Enum(common.Status(item.Status)),
	}
	return res
}

func (c *GoodsCategoryCase) ConvertToModel(item *admin.GoodsCategoryForm) *models.GoodsCategory {
	res := &models.GoodsCategory{
		ID:       item.GetId(),
		ParentID: item.GetParentId(),
		Picture:  item.GetPicture(),
		Name:     item.GetName(),
		Sort:     item.GetSort(),
		Status:   int32(item.GetStatus()),
	}
	return res
}
