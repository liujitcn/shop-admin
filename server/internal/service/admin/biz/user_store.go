package biz

import (
	"context"

	_string "github.com/liujitcn/go-utils/string"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	"github.com/liujitcn/shop-admin/server/internal/data"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type UserStoreCase struct {
	tx genData.Transaction
	data.UserStoreRepo
	baseAreaCase *BaseAreaCase
	baseUserCase *BaseUserCase
	baseRoleCase *BaseRoleCase
}

// NewUserStoreCase new a UserStore use case.
func NewUserStoreCase(tx genData.Transaction,
	userStoreRepo data.UserStoreRepo,
	baseAreaCase *BaseAreaCase,
	baseUserCase *BaseUserCase,
	baseRoleCase *BaseRoleCase,
) *UserStoreCase {
	return &UserStoreCase{
		tx:            tx,
		UserStoreRepo: userStoreRepo,
		baseAreaCase:  baseAreaCase,
		baseUserCase:  baseUserCase,
		baseRoleCase:  baseRoleCase,
	}
}

func (c *UserStoreCase) Page(ctx context.Context, req *admin.PageUserStoreRequest) (*admin.PageUserStoreResponse, error) {
	condition := &data.UserStoreCondition{
		Name:   req.GetName(),
		Status: int32(req.GetStatus()),
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}
	// 查询用户信息
	userIds := make([]int64, 0)
	for _, item := range page {
		userIds = append(userIds, item.UserID)
	}
	baseUserList := make([]*models.BaseUser, 0)
	baseUserList, err = c.baseUserCase.FindAll(ctx, &data.BaseUserCondition{
		Ids: userIds,
	})
	baseUserMap := make(map[int64]*models.BaseUser)
	for _, item := range baseUserList {
		baseUserMap[item.ID] = item
	}

	list := make([]*admin.UserStore, 0)
	for _, item := range page {
		userStore := c.ConvertToProto(ctx, item)
		if v, ok := baseUserMap[item.UserID]; ok {
			userStore.NickName = v.NickName
			userStore.Phone = v.Phone
		}
		list = append(list, userStore)
	}

	return &admin.PageUserStoreResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *UserStoreCase) GetUserStore(ctx context.Context, id int64) (*admin.UserStore, error) {
	userStore, err := c.Find(ctx, &data.UserStoreCondition{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	res := c.ConvertToProto(ctx, userStore)

	// 查询用户信息
	var baseUser *models.BaseUser
	baseUser, err = c.baseUserCase.GetFromID(ctx, id)
	if err != nil {
		return nil, err
	}
	res.NickName = baseUser.NickName
	res.Phone = baseUser.Phone
	return res, nil
}

func (c *UserStoreCase) AuditUserStore(ctx context.Context, req *admin.AuditUserStoreForm) error {
	// 查询当前审核信息
	userStore, err := c.Find(ctx, &data.UserStoreCondition{
		Id: req.GetId(),
	})
	if err != nil {
		return err
	}
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		err = c.UpdateByUserIdAndId(ctx, userStore.UserID, &models.UserStore{
			ID:     req.GetId(),
			Status: int32(req.GetStatus()),
			Remark: req.GetRemark(),
		})
		if err != nil {
			return err
		}
		// 审核通过， 修改用户角色为用户
		var code string
		if req.GetStatus() == common.UserStoreStatus_APPROVED {
			code = _const.BaseRoleCode_User
		} else {
			// 修改用户角色为游客
			code = _const.BaseRoleCode_Guest
		}

		var baseRole *models.BaseRole
		baseRole, err = c.baseRoleCase.Find(ctx, &data.BaseRoleCondition{
			Code: code,
		})
		if err != nil {
			return err
		}
		return c.baseUserCase.UpdateByID(ctx, &models.BaseUser{
			ID:     userStore.UserID,
			RoleID: baseRole.ID,
		})
	})
}

func (c *UserStoreCase) ConvertToProto(ctx context.Context, item *models.UserStore) *admin.UserStore {
	return &admin.UserStore{
		Id:              item.ID,
		Name:            item.Name,
		Address:         c.baseAreaCase.GetAddressListByCode(ctx, item.Address),
		Detail:          item.Detail,
		Picture:         _string.ConvertJsonStringToStringArray(item.Picture),
		BusinessLicense: _string.ConvertJsonStringToStringArray(item.BusinessLicense),
		Status:          common.UserStoreStatus(item.Status),
		Remark:          item.Remark,
	}
}
