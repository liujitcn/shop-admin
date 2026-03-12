package data

import (
	"context"
	"github.com/liujitcn/shop-admin/server/internal/data/dto"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"time"
)

type GoodsCondition struct {
	Id             int64   `query:"type:eq;column:id"`
	Ids            []int64 `query:"type:in;column:id"`
	Name           string  `query:"type:contains;column:name"`
	CategoryId     int64   `query:"type:eq;column:category_id"`
	CategoryPath   string
	Status         int32      `query:"type:eq;column:status"`
	StartCreatedAt *time.Time `query:"type:gte;column:created_at"` // 创建开始时间
	EndCreatedAt   *time.Time `query:"type:lte;column:created_at"` // 创建结束时间
}

type GoodsRepo interface {
	genRepo.BaseRepo[models.Goods, GoodsCondition]
	AddSaleNum(ctx context.Context, id, saleNum int64) error
	SubSaleNum(ctx context.Context, id, saleNum int64) error
	GoodsCategorySummary(ctx context.Context) ([]*dto.GoodsCategorySummary, error)
}

type goodsRepo struct {
	genRepo.BaseRepo[models.Goods, GoodsCondition]
	data *genData.Data
}

func NewGoodsRepo(data *genData.Data) GoodsRepo {
	base := genRepo.NewBaseRepo[models.Goods, GoodsCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).Goods.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).Goods.ID
		},
		func(entity *models.Goods) int64 {
			return entity.ID
		},
		new(models.Goods),
		100,
	)
	return &goodsRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *goodsRepo) AddSaleNum(ctx context.Context, id, saleNum int64) error {
	q := r.data.Query(ctx).Goods
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id)).Update(q.RealSaleNum, q.RealSaleNum.Add(saleNum))
	return err
}

func (r *goodsRepo) SubSaleNum(ctx context.Context, id, saleNum int64) error {
	q := r.data.Query(ctx).Goods
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id), q.RealSaleNum.Gte(saleNum)).Update(q.RealSaleNum, q.RealSaleNum.Sub(saleNum))
	return err
}

func (r *goodsRepo) GoodsCategorySummary(ctx context.Context) ([]*dto.GoodsCategorySummary, error) {
	category := r.data.Query(ctx).GoodsCategory
	m := r.data.Query(ctx).Goods
	q := m.WithContext(ctx)
	q = q.Join(category, category.ID.EqCol(m.CategoryID))

	results := make([]*dto.GoodsCategorySummary, 0)

	q.Select(category.ParentID.As("category_id"), m.ID.Count().As("goods_count")).Group(category.ParentID)
	err := q.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
