package biz

import (
	"github.com/liujitcn/shop-admin/server/internal/data"
)

type PayBillCase struct {
	data.PayBillRepo
}

// NewPayBillCase new a ShopPayBill use case.
func NewPayBillCase(
	payBillRepo data.PayBillRepo,
) *PayBillCase {
	return &PayBillCase{
		PayBillRepo: payBillRepo,
	}
}
