package data

import (
	"context"
	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type PayBillCondition struct {
	Id       int64  `search:"type:eq;column:id"`
	BillDate string `search:"type:eq;column:bill_date"`
	BillType string `search:"type:eq;column:bill_type"`
}

type PayBillRepo interface {
	baseRepo.BaseRepo[models.PayBill, PayBillCondition]
}

type payBillRepo struct {
	baseRepo.BaseRepo[models.PayBill, PayBillCondition]
	data *genData.Data
}

func NewPayBillRepo(data *genData.Data) PayBillRepo {
	base := baseRepo.NewBaseRepo[models.PayBill, PayBillCondition](
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
	)
	return &payBillRepo{
		BaseRepo: base,
		data:     data,
	}
}
