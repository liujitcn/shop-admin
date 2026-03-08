package biz

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/liujitcn/go-utils/str"
	"github.com/liujitcn/go-utils/timeutil"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type BaseMenuCase struct {
	tx genData.Transaction
	data.BaseMenuRepo
	casbinRuleCase *CasbinRuleCase
}

// NewBaseMenuCase new a BaseMenu use case.
func NewBaseMenuCase(
	tx genData.Transaction,
	baseMenuRepo data.BaseMenuRepo,
	casbinRuleCase *CasbinRuleCase,
) *BaseMenuCase {
	return &BaseMenuCase{
		tx:             tx,
		BaseMenuRepo:   baseMenuRepo,
		casbinRuleCase: casbinRuleCase,
	}
}

func (c *BaseMenuCase) GetFromID(ctx context.Context, id int64) (*models.BaseMenu, error) {
	return c.Find(ctx, &data.BaseMenuCondition{
		Id: id,
	})
}

func (c *BaseMenuCase) Tree(ctx context.Context, condition *data.BaseMenuCondition) (*admin.TreeBaseMenuResponse, error) {
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &admin.TreeBaseMenuResponse{
		List: c.buildTree(list, 0),
	}, nil
}

func (c *BaseMenuCase) Option(ctx context.Context, condition *data.BaseMenuCondition) (*common.TreeOptionResponse, error) {
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &common.TreeOptionResponse{
		List: c.buildOption(list, 0),
	}, nil
}

func (c *BaseMenuCase) List(ctx context.Context, condition *data.BaseMenuCondition) ([]*models.BaseMenu, error) {
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *BaseMenuCase) Update(ctx context.Context, req *admin.BaseMenuForm) error {
	baseMenu := c.ConvertToModel(req)
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err := c.UpdateByID(ctx, baseMenu)
		if err != nil {
			return err
		}
		// 判断apis 是否修改
		err = c.casbinRuleCase.RebuildCasbinRuleByMenuId(ctx, baseMenu.ID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *BaseMenuCase) Delete(ctx context.Context, id string) error {
	ids := str.ConvertStringToInt64Array(id)
	for _, item := range ids {
		// 查询下级
		count, err := c.Count(ctx, &data.BaseMenuCondition{
			ParentId: &item,
		})
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("删除菜单失败,下面有菜单")
		}
	}
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err := c.BaseMenuRepo.Delete(ctx, ids)
		if err != nil {
			return err
		}
		return c.casbinRuleCase.DeleteCasbinRuleByMenuIds(ctx, ids)
	})
}

// BuildRouteTree 构建菜单树状
func (c *BaseMenuCase) BuildRouteTree(menuList []*models.BaseMenu, parentId int64) []*admin.RouteItem {
	var menuRouters []*admin.RouteItem
	for _, menu := range menuList {
		if menu.ParentID == parentId {
			var menuRouter admin.RouteItem
			menuRouter.Path = &menu.Path
			menuRouter.Name = &menu.Name
			menuRouter.Component = &menu.Component
			menuRouter.Redirect = &menu.Redirect
			var meta admin.RouteMeta
			err := json.Unmarshal([]byte(menu.Meta), &meta)
			if err != nil {
				log.Error(err)
				continue
			}
			menuRouter.Meta = &meta
			menuRouter.Children = c.BuildRouteTree(menuList, menu.ID)
			menuRouters = append(menuRouters, &menuRouter)
		}
	}
	return menuRouters
}

// buildTree 构建菜单树状
func (c *BaseMenuCase) buildTree(menuList []*models.BaseMenu, parentId int64) []*admin.BaseMenu {
	var res []*admin.BaseMenu
	for _, item := range menuList {
		if item.ParentID == parentId {

			var meta admin.BaseMenuMeta
			err := json.Unmarshal([]byte(item.Meta), &meta)
			if err != nil {
				log.Error(err)
				continue
			}
			menu := &admin.BaseMenu{
				Id:        item.ID,
				ParentId:  item.ParentID,
				Type:      common.BaseMenuType(item.Type),
				Path:      item.Path,
				Name:      item.Name,
				Component: item.Component,
				Redirect:  item.Redirect,
				Meta:      &meta,
				Sort:      item.Sort,
				Status:    common.Status(item.Status),
				CreatedAt: timeutil.TimeToTimeString(item.CreatedAt),
				UpdatedAt: timeutil.TimeToTimeString(item.UpdatedAt),
			}
			menu.Children = c.buildTree(menuList, item.ID)
			res = append(res, menu)
		}
	}
	return res
}

// buildOptionTree 构建菜单树形选择
func (c *BaseMenuCase) buildOption(menuList []*models.BaseMenu, parentId int64) []*common.TreeOptionResponse_Option {
	var res []*common.TreeOptionResponse_Option
	for _, item := range menuList {
		if item.ParentID == parentId {
			var meta admin.RouteMeta
			err := json.Unmarshal([]byte(item.Meta), &meta)
			if err != nil {
				log.Error(err)
				continue
			}
			menu := &common.TreeOptionResponse_Option{
				Label: *meta.Title,
				Value: item.ID,
			}
			menu.Children = c.buildOption(menuList, item.ID)
			res = append(res, menu)
		}
	}
	return res
}

func (c *BaseMenuCase) ConvertToProto(item *models.BaseMenu) *admin.BaseMenuForm {
	var meta admin.BaseMenuMeta
	err := json.Unmarshal([]byte(item.Meta), &meta)
	if err != nil {
		log.Error(err)
	}
	res := &admin.BaseMenuForm{
		Id:        item.ID,
		ParentId:  trans.Int64(item.ParentID),
		Type:      trans.Enum(common.BaseMenuType(item.Type)),
		Path:      item.Path,
		Name:      item.Name,
		Component: item.Component,
		Redirect:  item.Redirect,
		Meta:      &meta,
		Apis:      str.ConvertJsonStringToStringArray(item.Apis),
		Sort:      item.Sort,
		Status:    trans.Enum(common.Status(item.Status)),
	}
	return res
}

func (c *BaseMenuCase) ConvertToModel(item *admin.BaseMenuForm) *models.BaseMenu {
	meta, err := json.Marshal(item.GetMeta())
	if err != nil {
		log.Error(err)
	}
	res := &models.BaseMenu{
		ID:        item.GetId(),
		ParentID:  item.GetParentId(),
		Type:      int32(item.GetType()),
		Path:      item.GetPath(),
		Name:      item.GetName(),
		Component: item.GetComponent(),
		Redirect:  item.GetRedirect(),
		Meta:      string(meta),
		Apis:      str.ConvertStringArrayToString(item.GetApis()),
		Sort:      item.GetSort(),
		Status:    int32(item.GetStatus()),
	}
	return res
}
