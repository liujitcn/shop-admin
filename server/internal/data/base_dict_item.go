package data

import (
	"context"

	genData "github.com/liujitcn/shop-gorm-gen/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	genRepo "github.com/liujitcn/shop-gorm-gen/repo"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type BaseDictItemCondition struct {
	Id      int64   `query:"type:eq;column:id"`
	DictId  int64   `query:"type:eq;column:dict_id"`
	DictIds []int64 `query:"type:in;column:dict_id"`
	Status  int32   `query:"type:eq;column:status"`
	Label   string  `query:"type:contains;column:label"`
	Value   string  `query:"type:contains;column:value"`
}

type BaseDictItemRepo interface {
	genRepo.BaseRepo[models.BaseDictItem, BaseDictItemCondition]
	FindLabelByCodeAndValue(ctx context.Context, code, value string) (string, error)
}

type baseDictItemRepo struct {
	genRepo.BaseRepo[models.BaseDictItem, BaseDictItemCondition]
	data *genData.Data
}

func NewBaseDictItemRepo(data *genData.Data) BaseDictItemRepo {
	base := genRepo.NewBaseRepo[models.BaseDictItem, BaseDictItemCondition](
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
		100,
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
