package biz

import (
	"context"

	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type BaseDeptCase struct {
	data.BaseDeptRepo
}

// NewBaseDeptCase new a BaseDept use case.
func NewBaseDeptCase(baseDeptRepo data.BaseDeptRepo) *BaseDeptCase {
	return &BaseDeptCase{
		BaseDeptRepo: baseDeptRepo,
	}
}

func (c *BaseDeptCase) GetFromID(ctx context.Context, id int64) (*models.BaseDept, error) {
	return c.Find(ctx, &data.BaseDeptCondition{
		Id: id,
	})
}

func (c *BaseDeptCase) Tree(ctx context.Context, condition *data.BaseDeptCondition) (*admin.TreeBaseDeptResponse, error) {
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &admin.TreeBaseDeptResponse{
		List: c.buildTree(list, 0),
	}, nil
}

func (c *BaseDeptCase) Option(ctx context.Context, condition *data.BaseDeptCondition) (*common.TreeOptionResponse, error) {
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &common.TreeOptionResponse{
		List: c.buildOption(list, 0),
	}, nil
}

// buildTree 构建部门树状
func (c *BaseDeptCase) buildTree(deptList []*models.BaseDept, parentId int64) []*admin.BaseDept {
	var res []*admin.BaseDept
	for _, item := range deptList {
		if item.ParentID == parentId {
			dept := &admin.BaseDept{
				Id:        item.ID,
				ParentId:  item.ParentID,
				Name:      item.Name,
				Sort:      item.Sort,
				Status:    common.Status(item.Status),
				Remark:    item.Remark,
				CreatedAt: timeutil.TimeToTimeString(item.CreatedAt),
				UpdatedAt: timeutil.TimeToTimeString(item.UpdatedAt),
			}
			dept.Children = c.buildTree(deptList, item.ID)
			res = append(res, dept)
		}
	}
	return res
}

// buildTree 构建部门树形选择
func (c *BaseDeptCase) buildOption(deptList []*models.BaseDept, parentId int64) []*common.TreeOptionResponse_Option {
	var res []*common.TreeOptionResponse_Option
	for _, item := range deptList {
		if item.ParentID == parentId {
			dept := &common.TreeOptionResponse_Option{
				Label: item.Name,
				Value: item.ID,
			}
			dept.Children = c.buildOption(deptList, item.ID)
			res = append(res, dept)
		}
	}
	return res
}

func (c *BaseDeptCase) ConvertToProto(item *models.BaseDept) *admin.BaseDeptForm {
	res := &admin.BaseDeptForm{
		Id:       item.ID,
		ParentId: trans.Int64(item.ParentID),
		Name:     item.Name,
		Sort:     item.Sort,
		Status:   trans.Enum(common.Status(item.Status)),
		Remark:   item.Remark,
	}
	return res
}

func (c *BaseDeptCase) ConvertToModel(item *admin.BaseDeptForm) *models.BaseDept {
	res := &models.BaseDept{
		ID:       item.GetId(),
		ParentID: item.GetParentId(),
		Name:     item.GetName(),
		Sort:     item.GetSort(),
		Status:   int32(item.GetStatus()),
		Remark:   item.GetRemark(),
	}
	return res
}
