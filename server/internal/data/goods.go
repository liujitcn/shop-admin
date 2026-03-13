package data

import (
	"context"
	"time"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	"github.com/liujitcn/shop-admin/server/internal/data/dto"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsCondition struct {
	Id             int64   `search:"type:eq;column:id"`
	Ids            []int64 `search:"type:in;column:id"`
	Name           string  `search:"type:contains;column:name"`
	CategoryId     int64   `search:"type:eq;column:category_id"`
	CategoryPath   string
	Status         int32      `search:"type:eq;column:status"`
	StartCreatedAt *time.Time `search:"type:gte;column:created_at"` // 创建开始时间
	EndCreatedAt   *time.Time `search:"type:lte;column:created_at"` // 创建结束时间
}

type GoodsRepo interface {
	baseRepo.BaseRepo[models.Goods, GoodsCondition]
	GoodsCategorySummary(ctx context.Context) ([]*dto.GoodsCategorySummary, error)
}

type goodsRepo struct {
	baseRepo.BaseRepo[models.Goods, GoodsCondition]
	data *genData.Data
}

func NewGoodsRepo(data *genData.Data) GoodsRepo {
	base := baseRepo.NewBaseRepo[models.Goods, GoodsCondition](
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
	)
	return &goodsRepo{
		BaseRepo: base,
		data:     data,
	}
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
