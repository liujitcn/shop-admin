package dto

// OrderGoodsStatusSummary 订单商品状态汇总
type OrderGoodsStatusSummary struct {
	Status     int32 `json:"status"`      // 状态
	GoodsCount int64 `json:"goods_count"` // 商品数量
	CategoryId int64 `json:"category_id"` // 分类
}

// OrderGoodsSummary 订单商品状态汇总
type OrderGoodsSummary struct {
	GoodsCount int64 `json:"goods_count"` // 商品数量
	GoodsId    int64 `json:"goods_id"`    // 商品id
}
