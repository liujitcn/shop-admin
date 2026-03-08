package data

import (
	"context"

	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
)

type ShopHotGoodsCondition struct {
	HotItemId int64
}

type ShopHotGoodsRepo interface {
	Delete(ctx context.Context, hotItemId int64) error
	Create(ctx context.Context, hotItemId int64, goodsIds []int64) error
	FindAll(ctx context.Context, hotItemId int64) ([]int64, error)
	ListPage(ctx context.Context, page, size int64, condition *ShopHotGoodsCondition) ([]int64, int64, error)
}

type shopHotGoodsRepo struct {
	data *genData.Data
}

func NewShopHotGoodsRepo(data *genData.Data) ShopHotGoodsRepo {
	return &shopHotGoodsRepo{data: data}
}

func (r *shopHotGoodsRepo) Delete(ctx context.Context, hotItemId int64) error {
	if hotItemId == 0 {
		return nil
	}
	q := r.data.Query(ctx).ShopHotGoods
	_, err := q.WithContext(ctx).Where(q.HotItemID.Eq(hotItemId)).Delete()
	return err
}

func (r *shopHotGoodsRepo) Create(ctx context.Context, hotItemId int64, goodsIds []int64) error {
	list := make([]*models.ShopHotGoods, 0)
	for idx, goodsId := range goodsIds {
		list = append(list, &models.ShopHotGoods{
			HotItemID: hotItemId,
			GoodsID:   goodsId,
			Sort:      int64(idx + 1),
		})
	}
	q := r.data.Query(ctx).ShopHotGoods
	err := q.WithContext(ctx).Clauses().CreateInBatches(list, 100)
	return err
}

func (r *shopHotGoodsRepo) FindAll(ctx context.Context, hotItemId int64) ([]int64, error) {
	m := r.data.Query(ctx).ShopHotGoods
	q := m.WithContext(ctx)
	q = q.Where(m.HotItemID.Eq(hotItemId))
	q = q.Order(m.Sort.Asc())
	list, err := q.Find()
	if err != nil {
		return nil, err
	}
	goodsIds := make([]int64, 0)
	for _, item := range list {
		goodsIds = append(goodsIds, item.GoodsID)
	}
	return goodsIds, nil
}
func (r *shopHotGoodsRepo) ListPage(ctx context.Context, page, size int64, condition *ShopHotGoodsCondition) ([]int64, int64, error) {
	m := r.data.Query(ctx).ShopHotGoods
	q := m.WithContext(ctx)
	if condition.HotItemId > 0 {
		q = q.Where(m.HotItemID.Eq(condition.HotItemId))
	}
	q = q.Order(m.Sort.Asc())
	offset, limit := genRepo.PageOffsetLimit(page, size)
	list, count, err := q.FindByPage(int(offset), int(limit))
	if err != nil {
		return nil, 0, err
	}
	goodsIds := make([]int64, 0)
	for _, item := range list {
		goodsIds = append(goodsIds, item.GoodsID)
	}
	return goodsIds, count, nil
}
