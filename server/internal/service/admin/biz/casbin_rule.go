package biz

import (
	"context"
	"slices"

	_string "github.com/liujitcn/go-utils/string"
	"github.com/liujitcn/kratos-kit/auth"
	authzEngine "github.com/liujitcn/kratos-kit/auth/authz/engine"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"

	"github.com/liujitcn/kratos-kit/auth/authz/engine/casbin"
)

type CasbinRuleCase struct {
	data.CasbinRuleRepo
	baseMenuRepo data.BaseMenuRepo
	baseRoleRepo data.BaseRoleRepo
	baseApiCase  *BaseApiCase
	authzEngine  authzEngine.Engine
}

// NewCasbinRuleCase new a CasbinRule use case.
func NewCasbinRuleCase(
	casbinRuleRepo data.CasbinRuleRepo,
	baseMenuRepo data.BaseMenuRepo,
	baseRoleRepo data.BaseRoleRepo,
	baseApiCase *BaseApiCase,
	authzEngine authzEngine.Engine,
) (*CasbinRuleCase, error) {
	c := &CasbinRuleCase{
		CasbinRuleRepo: casbinRuleRepo,
		baseMenuRepo:   baseMenuRepo,
		baseRoleRepo:   baseRoleRepo,
		baseApiCase:    baseApiCase,
		authzEngine:    authzEngine,
	}
	// 项目启动，加载casbin
	err := c.rebuildPolicyRule(context.Background())
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CasbinRuleCase) RebuildCasbinRuleByMenuId(ctx context.Context, menuId int64) error {
	// 查询全部角色
	baseRoleList, err := c.baseRoleRepo.FindAll(ctx, &data.BaseRoleCondition{})
	if err != nil {
		return err
	}
	for _, item := range baseRoleList {
		menus := _string.ConvertJsonStringToInt64Array(item.Menus)
		if slices.Contains(menus, menuId) {
			err = c.RebuildCasbinRuleByRole(ctx, item)
			if err != nil {
				return err
			}
		}
	}
	return c.rebuildPolicyRule(ctx)
}

func (c *CasbinRuleCase) DeleteCasbinRuleByMenuIds(ctx context.Context, menuIds []int64) error {
	// 查询全部角色
	baseRoleList, err := c.baseRoleRepo.FindAll(ctx, &data.BaseRoleCondition{})
	if err != nil {
		return err
	}
	for _, item := range baseRoleList {
		oldMenus := _string.ConvertJsonStringToInt64Array(item.Menus)
		newMenus := make([]int64, 0)
		for _, menuId := range oldMenus {
			// 不在删除列表
			if !slices.Contains(menuIds, menuId) {
				newMenus = append(newMenus, menuId)
			}
		}
		if len(oldMenus) != len(newMenus) {
			err = c.RebuildCasbinRuleByRole(ctx, item)
			if err != nil {
				return err
			}
		}
	}
	return c.rebuildPolicyRule(ctx)
}

func (c *CasbinRuleCase) RebuildCasbinRuleByRole(ctx context.Context, baseRole *models.BaseRole) error {
	// 删除casbin
	err := c.Delete(ctx, []string{baseRole.Code})
	if err != nil {
		return err
	}
	casbinRuleList := make([]*models.CasbinRule, 0)
	// 查询当前角色菜单

	// 查询菜单
	menuIds := _string.ConvertJsonStringToInt64Array(baseRole.Menus)
	if len(menuIds) == 0 {
		return nil
	}
	baseMenuList := make([]*models.BaseMenu, 0)
	baseMenuList, err = c.baseMenuRepo.FindAll(ctx, &data.BaseMenuCondition{
		Ids: menuIds,
	})

	operations := make([]string, 0)
	for _, item := range baseMenuList {
		apis := _string.ConvertJsonStringToStringArray(item.Apis)
		operations = append(operations, apis...)
	}
	if len(operations) == 0 {
		return nil
	}
	// 查询api列表
	baseApiList := make([]*models.BaseAPI, 0)
	allApiList := make([]*models.BaseAPI, 0)
	allApiList, err = c.baseApiCase.FindAll(ctx, &data.BaseApiCondition{})
	for _, item := range allApiList {
		if slices.Contains(operations, item.Operation) {
			baseApiList = append(baseApiList, item)
		}
	}
	for _, item := range baseApiList {
		casbinRuleList = append(casbinRuleList, &models.CasbinRule{
			Ptype: "p",
			V0:    baseRole.Code,
			V1:    item.Operation,
			V2:    string(auth.Action),
			V3:    "*",
		})
	}

	err = c.Create(ctx, casbinRuleList)
	if err != nil {
		return err
	}

	return c.rebuildPolicyRule(ctx)
}

func (c *CasbinRuleCase) DeleteCasbinRuleByRoleIds(ctx context.Context, roleIds []int64) error {
	baseRoleList, err := c.baseRoleRepo.FindAll(ctx, &data.BaseRoleCondition{
		Ids: roleIds,
	})
	if err != nil {
		return err
	}
	roleKeys := make([]string, 0)
	for _, item := range baseRoleList {
		roleKeys = append(roleKeys, item.Code)
	}
	// 删除casbin
	err = c.Delete(ctx, roleKeys)
	if err != nil {
		return err
	}
	return c.rebuildPolicyRule(ctx)
}

func (c *CasbinRuleCase) rebuildPolicyRule(ctx context.Context) error {
	policyRule := make([]casbin.PolicyRule, 0)
	// 查询全部api，默认给super 配置
	baseApiList, err := c.baseApiCase.FindAll(ctx, &data.BaseApiCondition{})
	if err != nil {
		return err
	}
	for _, item := range baseApiList {
		policyRule = append(policyRule, casbin.PolicyRule{
			PType: "p",
			V0:    _const.BaseRoleCode_Super,
			V1:    item.Operation,
			V2:    string(auth.Action),
			V3:    "*",
		})
	}
	// 查询casbin
	casbinRuleList := make([]*models.CasbinRule, 0)
	casbinRuleList, err = c.FindAll(ctx)
	for _, item := range casbinRuleList {
		policyRule = append(policyRule, casbin.PolicyRule{
			PType: item.Ptype,
			V0:    item.V0,
			V1:    item.V1,
			V2:    item.V2,
			V3:    item.V3,
			V4:    item.V4,
			V5:    item.V5,
		})
	}
	policyMap := make(authzEngine.PolicyMap)
	policyMap["policies"] = policyRule
	roleMap := make(authzEngine.RoleMap)
	err = c.authzEngine.SetPolicies(ctx, policyMap, roleMap)
	if err != nil {
		return err
	}
	return nil
}
