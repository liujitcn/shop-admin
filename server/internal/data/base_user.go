package data

import (
	"context"
	"time"

	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseUserCondition struct {
	Id             int64      `query:"type:eq;column:id"`
	Ids            []int64    `query:"type:in;column:id"`
	DeptId         int64      `query:"type:eq;column:dept_id"`
	Status         int32      `query:"type:eq;column:status"`
	DeptPath       string     // 用户部门路径
	UserName       string     `query:"type:contains;column:user_name"` // 用户账号
	NickName       string     `query:"type:contains;column:nick_name"` // 用户昵称
	Phone          string     `query:"type:contains;column:phone"`     // 手机号码
	Openid         string     `query:"type:contains;column:openid"`    // 手机号码
	Keyword        string     // 关键字
	StartCreatedAt *time.Time `query:"type:gte;column:created_at"` // 创建开始时间
	EndCreatedAt   *time.Time `query:"type:lte;column:created_at"` // 创建结束时间
}

type BaseUserRepo interface {
	genRepo.BaseRepo[models.BaseUser, BaseUserCondition]
}

type baseUserRepo struct {
	genRepo.BaseRepo[models.BaseUser, BaseUserCondition]
	data *genData.Data
}

func NewBaseUserRepo(data *genData.Data) BaseUserRepo {
	base := genRepo.NewBaseRepo[models.BaseUser, BaseUserCondition](
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
		100,
	)
	return &baseUserRepo{
		BaseRepo: base,
		data:     data,
	}
}
