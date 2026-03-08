package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type GoodsSkuCondition struct {
	Id       int64
	Ids      []int64
	GoodsId  int64
	SkuCode  string
	SkuCodes []string
	Status   int32
}

type GoodsSkuRepo interface {
	genRepo.BaseRepo[models.GoodsSku, GoodsSkuCondition]
	DeleteByGoodsId(ctx context.Context, goodsId int64) error
	AddSaleNum(ctx context.Context, skuCode string, saleNum int64) error
	SubSaleNum(ctx context.Context, skuCode string, saleNum int64) error
}

type goodsSkuRepo struct {
	genRepo.BaseRepo[models.GoodsSku, GoodsSkuCondition]
	data *genData.Data
}

func NewGoodsSkuRepo(data *genData.Data) GoodsSkuRepo {
	base := genRepo.NewBaseRepo[models.GoodsSku, GoodsSkuCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).GoodsSku.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).GoodsSku.ID
		},
		func(entity *models.GoodsSku) int64 {
			return entity.ID
		},
		new(models.GoodsSku),
		100,
	)
	return &goodsSkuRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *goodsSkuRepo) DeleteByGoodsId(ctx context.Context, goodsId int64) error {
	q := r.data.Query(ctx).GoodsSku
	_, err := q.WithContext(ctx).Where(q.GoodsID.Eq(goodsId)).Delete()
	return err
}

func (r *goodsSkuRepo) AddSaleNum(ctx context.Context, skuCode string, saleNum int64) error {
	q := r.data.Query(ctx).GoodsSku
	updates := map[string]interface{}{
		"real_sale_num": q.RealSaleNum.Add(saleNum),
		"inventory":     q.Inventory.Sub(saleNum),
	}
	_, err := q.WithContext(ctx).Where(q.SkuCode.Eq(skuCode), q.Inventory.Gte(saleNum)).Updates(updates)
	return err
}

func (r *goodsSkuRepo) SubSaleNum(ctx context.Context, skuCode string, saleNum int64) error {
	q := r.data.Query(ctx).GoodsSku
	updates := map[string]interface{}{
		"real_sale_num": q.RealSaleNum.Sub(saleNum),
		"inventory":     q.Inventory.Add(saleNum),
	}
	_, err := q.WithContext(ctx).Where(q.SkuCode.Eq(skuCode), q.RealSaleNum.Gte(saleNum)).Updates(updates)
	return err
}
