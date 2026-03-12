package biz

import (
	"context"

	_string "github.com/liujitcn/go-utils/string"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type GoodsSpecCase struct {
	data.GoodsSpecRepo
}

// NewGoodsSpecCase new a GoodsSpec use case.
func NewGoodsSpecCase(goodsSpecRepo data.GoodsSpecRepo) *GoodsSpecCase {
	return &GoodsSpecCase{
		GoodsSpecRepo: goodsSpecRepo,
	}
}

func (c *GoodsSpecCase) ListByGoodsId(ctx context.Context, goodsId int64) ([]*admin.GoodsSpec, error) {
	all, err := c.FindAll(ctx, &data.GoodsSpecCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*admin.GoodsSpec, 0)
	for _, item := range all {
		list = append(list, c.ConvertToProto(item))
	}
	return list, nil
}

func (c *GoodsSpecCase) BatchCreate(ctx context.Context, goodsId int64, spec []*admin.GoodsSpec) error {
	if len(spec) == 0 {
		return nil
	}
	// 查询旧数据
	oldSpecList, err := c.FindAll(ctx, &data.GoodsSpecCondition{
		GoodsId: goodsId,
	})
	if err != nil {
		return err
	}
	oldSpecIdMap := make(map[string]int64)
	for _, oldSpec := range oldSpecList {
		oldSpecIdMap[oldSpec.Name] = oldSpec.ID
	}
	specList := make([]*models.GoodsSpec, 0)
	for _, item := range spec {
		if id, ok := oldSpecIdMap[item.Name]; ok {
			item.Id = id
			err = c.UpdateByID(ctx, c.ConvertToModel(item))
			if err != nil {
				return err
			}
			delete(oldSpecIdMap, item.Name)
		} else {
			item.GoodsId = goodsId
			specList = append(specList, c.ConvertToModel(item))
		}
	}
	if len(oldSpecIdMap) > 0 {
		oldSpecId := make([]int64, 0)
		for _, v := range oldSpecIdMap {
			oldSpecId = append(oldSpecId, v)
		}
		err = c.Delete(ctx, oldSpecId)
		if err != nil {
			return err
		}
	}
	if len(specList) >= 0 {
		return c.GoodsSpecRepo.BatchCreate(ctx, specList)
	}
	return nil
}
func (c *GoodsSpecCase) ConvertToProto(item *models.GoodsSpec) *admin.GoodsSpec {
	res := &admin.GoodsSpec{
		Id:      item.ID,
		GoodsId: item.GoodsID,
		Name:    item.Name,
		Item:    _string.ConvertJsonStringToStringArray(item.Item),
		Sort:    item.Sort,
	}
	return res
}

func (c *GoodsSpecCase) ConvertToModel(item *admin.GoodsSpec) *models.GoodsSpec {
	res := &models.GoodsSpec{
		ID:      item.GetId(),
		GoodsID: item.GetGoodsId(),
		Name:    item.GetName(),
		Item:    _string.ConvertStringArrayToString(item.GetItem()),
		Sort:    item.GetSort(),
	}
	return res
}
