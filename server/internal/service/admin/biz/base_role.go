package biz

import (
	"context"
	"errors"

	_string "github.com/liujitcn/go-utils/string"
	_time "github.com/liujitcn/go-utils/time"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	"github.com/liujitcn/shop-admin/server/internal/data"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type BaseRoleCase struct {
	tx genData.Transaction
	data.BaseRoleRepo
	casbinRuleCase *CasbinRuleCase
}

// NewBaseRoleCase new a BaseRole use case.
func NewBaseRoleCase(
	tx genData.Transaction,
	baseRoleRepo data.BaseRoleRepo,
	casbinRuleCase *CasbinRuleCase,
) *BaseRoleCase {
	return &BaseRoleCase{
		tx:             tx,
		BaseRoleRepo:   baseRoleRepo,
		casbinRuleCase: casbinRuleCase,
	}
}

func (c *BaseRoleCase) GetFromID(ctx context.Context, id int64) (*models.BaseRole, error) {
	return c.Find(ctx, &data.BaseRoleCondition{Id: id})
}

func (c *BaseRoleCase) List(ctx context.Context, condition *data.BaseRoleCondition) ([]*models.BaseRole, error) {
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *BaseRoleCase) Page(ctx context.Context, req *admin.PageBaseRoleRequest) (*admin.PageBaseRoleResponse, error) {
	condition := &data.BaseRoleCondition{
		Status: int32(req.GetStatus()),
		Name:   req.GetName(),
		Code:   req.GetCode(),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.BaseRole, 0)
	for _, item := range page {
		list = append(list, &admin.BaseRole{
			Id:        item.ID,
			Name:      item.Name,
			Code:      item.Code,
			DataScope: common.BaseRoleDataScope(item.DataScope),
			Menus:     _string.ConvertJsonStringToInt64Array(item.Menus),
			Status:    common.Status(item.Status),
			Remark:    item.Remark,
			CreatedAt: _time.TimeToTimeString(item.CreatedAt),
			UpdatedAt: _time.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageBaseRoleResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseRoleCase) Create(ctx context.Context, req *admin.BaseRoleForm) error {
	baseRole := c.ConvertToModel(req)
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err := c.BaseRoleRepo.Create(ctx, baseRole)
		if err != nil {
			return err
		}
		return c.casbinRuleCase.RebuildCasbinRuleByRole(ctx, baseRole)
	})
}

func (c *BaseRoleCase) Update(ctx context.Context, req *admin.BaseRoleForm) error {
	oldBaseRole, err := c.GetFromID(ctx, req.GetId())
	if err != nil {
		return err
	}
	if oldBaseRole.Code == _const.BaseRoleCode_Super {
		return errors.New("更新角色失败，不能操作超级管理员角色")
	}
	baseRole := c.ConvertToModel(req)
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err = c.UpdateByID(ctx, baseRole)
		if err != nil {
			return err
		}
		err = c.casbinRuleCase.RebuildCasbinRuleByRole(ctx, baseRole)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *BaseRoleCase) Delete(ctx context.Context, id string) error {
	ids := _string.ConvertStringToInt64Array(id)
	count, err := c.Count(ctx, &data.BaseRoleCondition{
		Ids:  ids,
		Code: _const.BaseRoleCode_Super,
	})
	if err != nil {
		return errors.New("删除角色失败")
	}
	if count > 0 {
		return errors.New("删除角色失败，不能操作超级管理员角色")
	}

	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err = c.BaseRoleRepo.Delete(ctx, ids)
		if err != nil {
			return err
		}
		return c.casbinRuleCase.DeleteCasbinRuleByRoleIds(ctx, ids)
	})
}

func (c *BaseRoleCase) SetBaseRoleMenus(ctx context.Context, req *admin.SetMenusRequest) error {
	oldBaseRole, err := c.GetFromID(ctx, req.GetId())
	if err != nil {
		return err
	}
	if oldBaseRole.Code == _const.BaseRoleCode_Super {
		return errors.New("更新角色失败，不能操作超级管理员角色")
	}
	baseRole := &models.BaseRole{
		ID:    req.GetId(),
		Menus: _string.ConvertInt64ArrayToString(req.GetMenus()),
	}
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err = c.UpdateByID(ctx, baseRole)
		if err != nil {
			return err
		}
		baseRole.Code = oldBaseRole.Code
		err = c.casbinRuleCase.RebuildCasbinRuleByRole(ctx, baseRole)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *BaseRoleCase) ConvertToProto(item *models.BaseRole) *admin.BaseRoleForm {
	res := &admin.BaseRoleForm{
		Id:        item.ID,
		Name:      item.Name,
		Code:      item.Code,
		DataScope: trans.Enum(common.BaseRoleDataScope(item.DataScope)),
		Menus:     _string.ConvertJsonStringToInt64Array(item.Menus),
		Status:    trans.Enum(common.Status(item.Status)),
		Remark:    item.Remark,
	}
	return res
}

func (c *BaseRoleCase) ConvertToModel(item *admin.BaseRoleForm) *models.BaseRole {
	res := &models.BaseRole{
		ID:        item.GetId(),
		Name:      item.GetName(),
		Code:      item.GetCode(),
		DataScope: int32(item.GetDataScope()),
		Menus:     _string.ConvertInt64ArrayToString(item.GetMenus()),
		Status:    int32(item.GetStatus()),
		Remark:    item.GetRemark(),
	}
	return res
}
