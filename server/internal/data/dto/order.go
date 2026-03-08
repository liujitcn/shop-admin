package dto

// OrderSummary 月份订单汇总
type OrderSummary struct {
	Key        int64 `json:"key"` // 格式: "04"
	OrderCount int64 `json:"order_count"`
	SaleAmount int64 `json:"sale_amount"`
}
