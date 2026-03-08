package data

import (
	"context"
	"errors"

	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type UserStoreCondition struct {
	Id     int64
	Name   string
	UserId int64
	Status int32
}

type UserStoreRepo interface {
	genRepo.BaseRepo[models.UserStore, UserStoreCondition]
	DeleteByUserIdAndIds(ctx context.Context, userID int64, ids []int64) error
	UpdateByUserIdAndId(ctx context.Context, userID int64, userStore *models.UserStore) error
}

type userStoreRepo struct {
	genRepo.BaseRepo[models.UserStore, UserStoreCondition]
	data *genData.Data
}

func NewUserStoreRepo(data *genData.Data) UserStoreRepo {
	base := genRepo.NewBaseRepo[models.UserStore, UserStoreCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).UserStore.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).UserStore.ID
		},
		func(entity *models.UserStore) int64 {
			return entity.ID
		},
		new(models.UserStore),
		100,
	)
	return &userStoreRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *userStoreRepo) DeleteByUserIdAndIds(ctx context.Context, userID int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	q := r.data.Query(ctx).UserStore
	_, err := q.WithContext(ctx).Where(q.UserID.Eq(userID), q.ID.In(ids...)).Delete()
	return err
}

func (r *userStoreRepo) UpdateByUserIdAndId(ctx context.Context, userID int64, userStore *models.UserStore) error {
	if userStore.ID == 0 {
		return errors.New("userStore can not update without id")
	}
	m := r.data.Query(ctx).UserStore
	q := m.WithContext(ctx)
	if userID > 0 {
		q = q.Where(m.UserID.Eq(userID))
	}
	_, err := q.Updates(userStore)
	return err
}
