package data

import (
	"context"

	baseRepo "github.com/liujitcn/gorm-kit/repo"
	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseDictItemCondition struct {
	Id      int64   `search:"type:eq;column:id"`
	DictId  int64   `search:"type:eq;column:dict_id"`
	DictIds []int64 `search:"type:in;column:dict_id"`
	Status  int32   `search:"type:eq;column:status"`
	Label   string  `search:"type:contains;column:label"`
	Value   string  `search:"type:contains;column:value"`
}

type BaseDictItemRepo interface {
	baseRepo.BaseRepo[models.BaseDictItem, BaseDictItemCondition]
	FindLabelByCodeAndValue(ctx context.Context, code, value string) (string, error)
}

type baseDictItemRepo struct {
	baseRepo.BaseRepo[models.BaseDictItem, BaseDictItemCondition]
	data *genData.Data
}

func NewBaseDictItemRepo(data *genData.Data) BaseDictItemRepo {
	base := baseRepo.NewBaseRepo[models.BaseDictItem, BaseDictItemCondition](
		func(ctx context.Context) gen.Dao {
			return new(data.Query(ctx).BaseDictItem.WithContext(ctx).DO)
		},
		func(ctx context.Context) field.Int64 {
			return data.Query(ctx).BaseDictItem.ID
		},
		func(entity *models.BaseDictItem) int64 {
			return entity.ID
		},
		new(models.BaseDictItem),
	)
	return &baseDictItemRepo{
		BaseRepo: base,
		data:     data,
	}
}

func (r *baseDictItemRepo) FindLabelByCodeAndValue(ctx context.Context, code, value string) (string, error) {
	m := r.data.Query(ctx).BaseDictItem
	dict := r.data.Query(ctx).BaseDict
	q := m.WithContext(ctx).Select(m.Label).Join(dict, m.DictID.EqCol(dict.ID))
	q = q.Where(dict.Code.Eq(code))
	q = q.Where(m.Value.Eq(value))

	var label string
	err := q.Scan(&label)
	return label, err
}
