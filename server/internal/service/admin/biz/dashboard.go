package biz

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/liujitcn/go-utils"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-admin/server/internal/data/dto"
	"github.com/liujitcn/shop-gorm-gen/models"
)

type DashboardCase struct {
	baseUserCase      *BaseUserCase
	goodsCase         *GoodsCase
	goodsCategoryCase *GoodsCategoryCase
	orderCase         *OrderCase
	orderGoodsCase    *OrderGoodsCase
	baseDictCase      *BaseDictCase
	baseDictItemCase  *BaseDictItemCase
}

// NewDashboardCase new a Dashboard use case.
func NewDashboardCase(
	baseUserCase *BaseUserCase,
	goodsCase *GoodsCase,
	goodsCategoryCase *GoodsCategoryCase,
	orderCase *OrderCase,
	orderGoodsCase *OrderGoodsCase,
	baseDictCase *BaseDictCase,
	baseDictItemCase *BaseDictItemCase,
) *DashboardCase {
	return &DashboardCase{
		baseUserCase:      baseUserCase,
		goodsCase:         goodsCase,
		goodsCategoryCase: goodsCategoryCase,
		orderCase:         orderCase,
		orderGoodsCase:    orderGoodsCase,
		baseDictCase:      baseDictCase,
		baseDictItemCase:  baseDictItemCase,
	}
}

func (c *DashboardCase) DashboardCountUser(ctx context.Context, req *admin.DashboardCountRequest) (*admin.DashboardCountResponse, error) {
	var err error
	startCreatedAt, endCreatedAt := utils.GetCreatedAt(utils.TimeType(req.GetTimeType()))

	// 查询今日注册用户数
	var user admin.DashboardCountResponse
	user.NewNum, err = c.baseUserCase.Count(ctx, &data.BaseUserCondition{
		StartCreatedAt: &startCreatedAt,
		EndCreatedAt:   &endCreatedAt,
	})
	if err != nil {
		return nil, err
	}
	// 查询总用户数
	user.TotalNum, err = c.baseUserCase.Count(ctx, &data.BaseUserCondition{})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *DashboardCase) DashboardCountGoods(ctx context.Context, req *admin.DashboardCountRequest) (*admin.DashboardCountResponse, error) {
	var err error
	startCreatedAt, endCreatedAt := utils.GetCreatedAt(utils.TimeType(req.GetTimeType()))

	// 查询今日新增商品
	var goods admin.DashboardCountResponse
	goods.NewNum, err = c.goodsCase.Count(ctx, &data.GoodsCondition{
		StartCreatedAt: &startCreatedAt,
		EndCreatedAt:   &endCreatedAt,
	})
	if err != nil {
		return nil, err
	}
	// 查询总商品数
	goods.TotalNum, err = c.goodsCase.Count(ctx, &data.GoodsCondition{})
	if err != nil {
		return nil, err
	}

	return &goods, nil
}

func (c *DashboardCase) DashboardCountOrder(ctx context.Context, req *admin.DashboardCountRequest) (*admin.DashboardCountResponse, error) {
	var err error
	startCreatedAt, endCreatedAt := utils.GetCreatedAt(utils.TimeType(req.GetTimeType()))

	// 查询今日订单数
	var order admin.DashboardCountResponse
	order.NewNum, err = c.orderCase.Count(ctx, &data.OrderCondition{
		StartCreatedAt: &startCreatedAt,
		EndCreatedAt:   &endCreatedAt,
	})
	if err != nil {
		return nil, err
	}
	// 查询总订单数
	order.TotalNum, err = c.orderCase.Count(ctx, &data.OrderCondition{})
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (c *DashboardCase) DashboardCountSale(ctx context.Context, req *admin.DashboardCountRequest) (*admin.DashboardCountResponse, error) {
	var err error
	startCreatedAt, endCreatedAt := utils.GetCreatedAt(utils.TimeType(req.GetTimeType()))

	// 查询今日销售额
	var sale admin.DashboardCountResponse
	sale.NewNum, err = c.orderCase.Sum(ctx, &data.OrderCondition{
		StartCreatedAt: &startCreatedAt,
		EndCreatedAt:   &endCreatedAt,
	})
	if err != nil {
		return nil, err
	}
	sale.TotalNum, err = c.orderCase.Sum(ctx, &data.OrderCondition{})
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (c *DashboardCase) DashboardBarOrder(ctx context.Context, req *admin.DashboardBarOrderRequest) (*admin.DashboardBarResponse, error) {
	var xAxisInt int
	switch req.GetTimeType() {
	case admin.DashboardTimeType_MONTH:
		xAxisInt = 12
	case admin.DashboardTimeType_DAY:
		xAxisInt = 7
	default:
		now := time.Now()
		year, month, _ := now.Date()
		startCreatedAt := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		endCreatedAt := startCreatedAt.AddDate(0, 1, 0).Add(-time.Nanosecond)
		xAxisInt = endCreatedAt.Day()
	}
	startCreatedAt, endCreatedAt := utils.GetCreatedAt(utils.TimeType(req.GetTimeType()))
	summary, err := c.orderCase.OrderSummary(ctx, int32(req.GetTimeType()), &data.OrderCondition{
		StartCreatedAt: &startCreatedAt,
		EndCreatedAt:   &endCreatedAt,
	})
	if err != nil {
		return nil, err
	}
	summaryMap := make(map[int64]*dto.OrderSummary)
	for _, item := range summary {
		// 周统计
		if req.GetTimeType() == admin.DashboardTimeType_WEEK {
			// 周日转换成7
			if item.Key == 1 {
				item.Key = 7
			} else {
				item.Key -= 1
			}
		}
		summaryMap[item.Key] = item
	}

	axisData := make([]string, 0)
	// 订单数量
	orderCountRow := make([]int64, 0)
	// 销售金额
	saleAmountRow := make([]int64, 0)
	// 订单数量增长率
	orderCountRateRow := make([]int64, 0)
	// 销售金额增长率
	saleAmountRateRow := make([]int64, 0)
	for i := 0; i < xAxisInt; i++ {
		month := int64(i + 1)
		axisData = append(axisData, utils.FormatDate(utils.TimeType(req.GetTimeType()), i))

		if item, ok := summaryMap[month]; ok {
			orderCountRow = append(orderCountRow, item.OrderCount)
			saleAmountRow = append(saleAmountRow, item.SaleAmount)
		} else {
			orderCountRow = append(orderCountRow, 0)
			saleAmountRow = append(saleAmountRow, 0)
		}
		if i == 0 {
			orderCountRateRow = append(orderCountRateRow, utils.CalcGrowthRate(0, orderCountRow[i]))
			saleAmountRateRow = append(saleAmountRateRow, utils.CalcGrowthRate(0, saleAmountRow[i]))
		} else {
			orderCountRateRow = append(orderCountRateRow, utils.CalcGrowthRate(orderCountRow[i-1], orderCountRow[i]))
			saleAmountRateRow = append(saleAmountRateRow, utils.CalcGrowthRate(saleAmountRow[i-1], saleAmountRow[i]))
		}
	}

	seriesData := make([]*admin.DashboardBarResponse_SeriesData, 0)
	seriesData = append(seriesData, &admin.DashboardBarResponse_SeriesData{
		Value: orderCountRow,
	})
	seriesData = append(seriesData, &admin.DashboardBarResponse_SeriesData{
		Value: saleAmountRow,
	})
	seriesData = append(seriesData, &admin.DashboardBarResponse_SeriesData{
		Value: orderCountRateRow,
	})
	seriesData = append(seriesData, &admin.DashboardBarResponse_SeriesData{
		Value: saleAmountRateRow,
	})

	return &admin.DashboardBarResponse{
		AxisData:   axisData,
		SeriesData: seriesData,
	}, nil
}

func (c *DashboardCase) DashboardBarGoods(ctx context.Context, req *admin.DashboardBarGoodsRequest) (*admin.DashboardBarResponse, error) {
	startCreatedAt, endCreatedAt := utils.GetCreatedAt(utils.TimeType(req.GetTimeType()))
	summary, err := c.orderGoodsCase.OrderGoodsSummary(ctx, req.GetTop(), &startCreatedAt, &endCreatedAt)
	if err != nil {
		return nil, err
	}
	goodsIds := make([]int64, 0)
	for _, item := range summary {
		goodsIds = append(goodsIds, item.GoodsId)
	}
	// 查询商品信息
	goodsList := make([]*models.Goods, 0)
	goodsList, err = c.goodsCase.FindAll(ctx, &data.GoodsCondition{
		Ids: goodsIds,
	})
	if err != nil {
		return nil, err
	}
	goodsNameMap := make(map[int64]string)
	for _, item := range goodsList {
		goodsNameMap[item.ID] = item.Name
	}

	// 换成正序
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].GoodsCount > summary[j].GoodsCount
	})

	axisData := make([]string, 0)
	// 商品数量
	goodsCountRow := make([]int64, 0)
	for _, item := range summary {
		axisData = append(axisData, goodsNameMap[item.GoodsId])
		goodsCountRow = append(goodsCountRow, item.GoodsCount)
	}
	seriesData := make([]*admin.DashboardBarResponse_SeriesData, 0)
	seriesData = append(seriesData, &admin.DashboardBarResponse_SeriesData{
		Value: goodsCountRow,
	})
	return &admin.DashboardBarResponse{
		AxisData:   axisData,
		SeriesData: seriesData,
	}, nil
}

func (c *DashboardCase) DashboardPieGoods(ctx context.Context, req *admin.DashboardPieGoodsRequest) (*admin.DashboardPieResponse, error) {
	summary, err := c.goodsCase.GoodsCategorySummary(ctx)
	if err != nil {
		return nil, err
	}
	// 查询一级分类
	var parentId int64 = 0
	goodsCategoryNameMap := c.goodsCategoryCase.NameMap(ctx, &data.GoodsCategoryCondition{
		ParentId: &parentId,
	})
	seriesData := make([]*admin.DashboardPieResponse_SeriesData, 0)
	for _, item := range summary {
		seriesData = append(seriesData, &admin.DashboardPieResponse_SeriesData{
			Value: item.GoodsCount,
			Name:  goodsCategoryNameMap[item.CategoryId],
		})
	}
	return &admin.DashboardPieResponse{
		SeriesData: seriesData,
	}, nil
}

func (c *DashboardCase) DashboardRadarOrder(ctx context.Context, req *admin.DashboardRadarOrderRequest) (*admin.DashboardRadarResponse, error) {
	startCreatedAt, endCreatedAt := utils.GetCreatedAt(utils.TimeType(req.GetTimeType()))
	summary, err := c.orderGoodsCase.OrderGoodsStatusSummary(ctx, &startCreatedAt, &endCreatedAt)
	if err != nil {
		return nil, err
	}

	summaryMap := make(map[string]int64)
	for _, item := range summary {
		summaryMap[fmt.Sprintf("%d_%d", item.CategoryId, item.Status)] = item.GoodsCount
	}

	// 查询一级分类
	var parentId int64 = 0
	goodsCategoryNameMap := c.goodsCategoryCase.NameMap(ctx, &data.GoodsCategoryCondition{
		ParentId: &parentId,
	})
	// 查询字典数据
	var baseDict *models.BaseDict
	baseDict, err = c.baseDictCase.Find(ctx, &data.BaseDictCondition{
		Code: "order_status",
	})
	if err != nil {
		return nil, err
	}
	baseDictItemList := make([]*models.BaseDictItem, 0)
	baseDictItemList, err = c.baseDictItemCase.FindAll(ctx, &data.BaseDictItemCondition{
		DictId: baseDict.ID,
	})
	// 数据不完整，直接返回
	if len(baseDictItemList) == 0 || len(goodsCategoryNameMap) == 0 {
		return &admin.DashboardRadarResponse{}, nil
	}
	legendData := make([]string, 0)
	seriesData := make([]*admin.DashboardRadarResponse_SeriesData, 0)
	radarIndicator := make([]*admin.DashboardRadarResponse_RadarIndicator, 0)
	for index, item := range baseDictItemList {
		legendData = append(legendData, item.Label)

		goodsNum := make([]int64, 0)
		for k, v := range goodsCategoryNameMap {
			if index == 0 {
				radarIndicator = append(radarIndicator, &admin.DashboardRadarResponse_RadarIndicator{
					Name: v,
				})
			}
			// 判断是非存在
			key := fmt.Sprintf("%d_%s", k, item.Value)
			if v, ok := summaryMap[key]; ok {
				goodsNum = append(goodsNum, v)
			} else {
				goodsNum = append(goodsNum, 0)
			}
		}
		seriesData = append(seriesData, &admin.DashboardRadarResponse_SeriesData{
			Name:  item.Label,
			Value: goodsNum,
		})
	}

	return &admin.DashboardRadarResponse{
		LegendData:     legendData,
		RadarIndicator: radarIndicator,
		SeriesData:     seriesData,
	}, nil
}
