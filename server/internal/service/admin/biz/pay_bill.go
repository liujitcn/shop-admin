package biz

import (
	"context"

	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type PayBillCase struct {
	data.PayBillRepo
}

// NewPayBillCase new a PayBill use case.
func NewPayBillCase(payBillRepo data.PayBillRepo) *PayBillCase {
	return &PayBillCase{
		PayBillRepo: payBillRepo,
	}
}

func (c *PayBillCase) GetFromID(ctx context.Context, id int64) (*models.PayBill, error) {
	return c.Find(ctx, &data.PayBillCondition{
		Id: id,
	})
}

func (c *PayBillCase) Page(ctx context.Context, req *admin.PagePayBillRequest) (*admin.PagePayBillResponse, error) {
	condition := &data.PayBillCondition{
		BillDate: req.GetBillDate(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.PayBill, 0)
	for _, item := range page {
		list = append(list, c.ConvertToProto(item))
	}

	return &admin.PagePayBillResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *PayBillCase) ConvertToProto(item *models.PayBill) *admin.PayBill {
	return &admin.PayBill{
		Id:               item.ID,
		BillDate:         item.BillDate,
		BillType:         item.BillType,
		FilePath:         item.FilePath,
		HashType:         item.HashType,
		HashValue:        item.HashValue,
		TotalCount:       item.TotalCount,
		TotalAmount:      item.TotalAmount,
		ThirdTotalCount:  item.ThirdTotalCount,
		ThirdTotalAmount: item.ThirdTotalAmount,
		Status:           common.PayBillStatus(item.Status),
	}
}
