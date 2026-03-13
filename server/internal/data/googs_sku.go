package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsSkuCondition struct {
	Id       int64    `search:"type:eq;column:id"`
	Ids      []int64  `search:"type:in;column:id"`
	GoodsId  int64    `search:"type:eq;column:goods_id"`
	SkuCode  string   `search:"type:eq;column:sku_code"`
	SkuCodes []string `search:"type:in;column:sku_code"`
	Status   int32    `search:"type:eq;column:status"`
}

type GoodsSkuRepo interface {
	baseRepo.BaseRepo[models.GoodsSKU, GoodsSkuCondition]
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
	AddSaleNum(ctx context.Context, skuCode string, saleNum int64) error
	SubSaleNum(ctx context.Context, skuCode string, saleNum int64) error
}

type goodsSkuRepo struct {
	baseRepo.BaseRepo[models.GoodsSKU, GoodsSkuCondition]
	data *genData.Data
}

func NewGoodsSkuRepo(data *genData.Data) GoodsSkuRepo {
	base := baseRepo.NewBaseRepo[models.GoodsSKU, GoodsSkuCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).GoodsSKU.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsSKU.ID
		},
		func(entity *models.GoodsSKU) int64 {
			return entity.ID
		},
		new(models.GoodsSKU),
	)
	return &goodsSkuRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *goodsSkuRepo) DeleteByGoodsId(ctx context.Context, goodsId int64) error {
	q := r.data.Query(ctx).GoodsSKU
	_, err := q.WithContext(ctx).Where(q.GoodsID.Eq(goodsId)).Delete()
	return err
}

func (r *goodsSkuRepo) AddSaleNum(ctx context.Context, skuCode string, saleNum int64) error {
	q := r.data.Query(ctx).GoodsSKU
	updates := map[string]interface{}{
		"real_sale_num": q.RealSaleNum.Add(saleNum),
		"inventory":     q.Inventory.Sub(saleNum),
	}
	_, err := q.WithContext(ctx).Where(q.SkuCode.Eq(skuCode), q.Inventory.Gte(saleNum)).Updates(updates)
	return err
}

func (r *goodsSkuRepo) SubSaleNum(ctx context.Context, skuCode string, saleNum int64) error {
	q := r.data.Query(ctx).GoodsSKU
	updates := map[string]interface{}{
		"real_sale_num": q.RealSaleNum.Sub(saleNum),
		"inventory":     q.Inventory.Add(saleNum),
	}
	_, err := q.WithContext(ctx).Where(q.SkuCode.Eq(skuCode), q.RealSaleNum.Gte(saleNum)).Updates(updates)
	return err
}
