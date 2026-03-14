package biz

import (
	"context"
	"fmt"

	"github.com/liujitcn/go-utils/crypto"
	_time "github.com/liujitcn/go-utils/time"
	"github.com/liujitcn/go-utils/trans"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/api/gen/go/common"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type BaseUserCase struct {
	data.BaseUserRepo
	baseDeptRepo data.BaseDeptRepo
}

// NewBaseUserCase new a BaseUser use case.
func NewBaseUserCase(baseUserRepo data.BaseUserRepo, deptRepo data.BaseDeptRepo) *BaseUserCase {
	return &BaseUserCase{
		BaseUserRepo: baseUserRepo,
		baseDeptRepo: deptRepo,
	}
}

func (c *BaseUserCase) GetFromID(ctx context.Context, id int64) (*models.BaseUser, error) {
	return c.Find(ctx, &data.BaseUserCondition{Id: id})
}

func (c *BaseUserCase) List(ctx context.Context) ([]*models.BaseUser, error) {
	condition := &data.BaseUserCondition{}
	list, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *BaseUserCase) Page(ctx context.Context, req *admin.PageBaseUserRequest) (*admin.PageBaseUserResponse, error) {
	condition := &data.BaseUserCondition{
		DeptId:   req.GetDeptId(),
		Status:   int32(req.GetStatus()),
		UserName: req.GetUserName(),
		NickName: req.GetNickName(),
		Phone:    req.GetPhone(),
	}
	if condition.DeptId > 0 {
		// 查询部门路径
		dept, err := c.baseDeptRepo.Find(ctx, &data.BaseDeptCondition{Id: condition.DeptId})
		if err != nil {
			return nil, err
		}
		condition.DeptId = 0
		condition.DeptPath = fmt.Sprintf("%s/", dept.Path)
	}
	page, count, err := c.ListPage(ctx, req.GetPageNum(), req.GetPageSize(), condition)
	if err != nil {
		return nil, err
	}

	list := make([]*admin.BaseUser, 0)
	for _, item := range page {
		list = append(list, &admin.BaseUser{
			Id:        item.ID,
			UserName:  item.UserName,
			NickName:  item.NickName,
			RoleId:    item.RoleID,
			DeptId:    item.DeptID,
			Phone:     item.Phone,
			Gender:    common.BaseUserGender(item.Gender),
			Avatar:    item.Avatar,
			Status:    common.Status(item.Status),
			Remark:    item.Remark,
			CreatedAt: _time.TimeToTimeString(item.CreatedAt),
			UpdatedAt: _time.TimeToTimeString(item.UpdatedAt),
		})
	}

	return &admin.PageBaseUserResponse{
		List:  list,
		Total: int32(count),
	}, nil
}

func (c *BaseUserCase) ConvertToProto(item *models.BaseUser) *admin.BaseUserForm {
	res := &admin.BaseUserForm{
		Id:       item.ID,
		UserName: item.UserName,
		NickName: item.NickName,
		RoleId:   trans.Int64(item.RoleID),
		DeptId:   trans.Int64(item.DeptID),
		Phone:    item.Phone,
		Gender:   trans.Enum(common.BaseUserGender(item.Gender)),
		Avatar:   item.Avatar,
		Status:   trans.Enum(common.Status(item.Status)),
		Remark:   item.Remark,
	}
	return res
}

func (c *BaseUserCase) ConvertToModel(item *admin.BaseUserForm, password string) *models.BaseUser {
	res := &models.BaseUser{
		ID:       item.GetId(),
		UserName: item.GetUserName(),
		NickName: item.GetNickName(),
		RoleID:   item.GetRoleId(),
		DeptID:   item.GetDeptId(),
		Phone:    item.GetPhone(),
		Password: password,
		Gender:   int32(item.GetGender()),
		Avatar:   item.GetAvatar(),
		Status:   int32(item.GetStatus()),
		Remark:   item.GetRemark(),
	}
	return res
}

func (c *BaseUserCase) GetDefaultPassword(userName, phone string) string {
	// 截取前四位，不够则取全部
	prefix := phone
	if len(phone) > 4 {
		prefix = phone[:4]
	}
	// 补0到四位
	prefix = fmt.Sprintf("%-4s", prefix)
	return fmt.Sprintf("%s@%s", userName, prefix)
}

func (c *BaseUserCase) GetPassword(passwordStr string) string {
	password, err := crypto.Encrypt(passwordStr)
	if err != nil {
		return ""
	}
	return password
}
