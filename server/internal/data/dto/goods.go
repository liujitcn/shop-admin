package dto

// GoodsCategorySummary 商品分类汇总
type GoodsCategorySummary struct {
	GoodsCount int64 `json:"goods_count"`
	CategoryId int64 `json:"category_id"`
}
