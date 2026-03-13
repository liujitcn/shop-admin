package data

import (
	"context"
	"time"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseUserCondition struct {
	Id             int64      `search:"type:eq;column:id"`
	Ids            []int64    `search:"type:in;column:id"`
	DeptId         int64      `search:"type:eq;column:dept_id"`
	Status         int32      `search:"type:eq;column:status"`
	DeptPath       string     // 用户部门路径
	UserName       string     `search:"type:eq;column:user_name"`       // 用户账号
	NickName       string     `search:"type:contains;column:nick_name"` // 用户昵称
	Phone          string     `search:"type:contains;column:phone"`     // 手机号码
	Openid         string     `search:"type:contains;column:openid"`    // 手机号码
	Keyword        string     // 关键字
	StartCreatedAt *time.Time `search:"type:gte;column:created_at"` // 创建开始时间
	EndCreatedAt   *time.Time `search:"type:lte;column:created_at"` // 创建结束时间
}

type BaseUserRepo interface {
	baseRepo.BaseRepo[models.BaseUser, BaseUserCondition]
}

type baseUserRepo struct {
	baseRepo.BaseRepo[models.BaseUser, BaseUserCondition]
	data *genData.Data
}

func NewBaseUserRepo(data *genData.Data) BaseUserRepo {
	base := baseRepo.NewBaseRepo[models.BaseUser, BaseUserCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseUser.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseUser.ID
		},
		func(entity *models.BaseUser) int64 {
			return entity.ID
		},
		new(models.BaseUser),
	)
	return &baseUserRepo{
		BaseRepo: base,
		data:     data,
	}
}
