package data

import (
	"context"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type PayBillCondition struct {
	Id       int64  `query:"type:eq;column:id"`
	BillDate string `query:"type:eq;column:bill_date"`
	BillType string `query:"type:eq;column:bill_type"`
}

type PayBillRepo interface {
	genRepo.BaseRepo[models.PayBill, PayBillCondition]
}

type payBillRepo struct {
	genRepo.BaseRepo[models.PayBill, PayBillCondition]
	data *genData.Data
}

func NewPayBillRepo(data *genData.Data) PayBillRepo {
	base := genRepo.NewBaseRepo[models.PayBill, PayBillCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).PayBill.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).PayBill.ID
		},
		func(entity *models.PayBill) int64 {
			return entity.ID
		},
		new(models.PayBill),
		100,
	)
	return &payBillRepo{
		BaseRepo: base,
		data:     data,
	}
}
