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
	Id             int64
	Ids            []int64
	DeptId         int64
	Status         int32
	DeptPath       string     // 用户部门路径
	UserName       string     // 用户账号
	NickName       string     // 用户昵称
	Phone          string     // 手机号码
	Openid         string     // 手机号码
	Keyword        string     // 关键字
	StartCreatedAt *time.Time // 创建开始时间
	EndCreatedAt   *time.Time // 创建结束时间
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
