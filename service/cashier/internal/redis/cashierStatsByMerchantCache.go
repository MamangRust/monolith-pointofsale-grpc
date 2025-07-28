package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

const (
	cashierStatsMonthTotalSalesByMerchantCacheKey = "cashier:stats:month:%d:year:%d:id:%d"
	cashierStatsYearTotalSalesByMerchantCacheKey  = "cashier:stats:year:%d:merchant:%d"

	cashierStatsMonthSalesByMerchantCacheKey = "cashier:stats:month:%d:merchant:%d"
	cashierStatsYearSalesByMerchantCacheKey  = "cashier:stats:year:%d:merchant:%d"
)

type cashierStatsByMerchantCache struct {
	store *CacheStore
}

func NewCashierStatsByMerchantCache(store *CacheStore) *cashierStatsByMerchantCache {
	return &cashierStatsByMerchantCache{store: store}
}

func (s *cashierStatsByMerchantCache) GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*response.CashierResponseMonthTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByMerchantCacheKey, req.Month, req.Year, req.MerchantID)
	result, found := GetFromCache[[]*response.CashierResponseMonthTotalSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant, res []*response.CashierResponseMonthTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByMerchantCacheKey, req.Month, req.Year, req.MerchantID)
	SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*response.CashierResponseYearTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := GetFromCache[[]*response.CashierResponseYearTotalSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant, res []*response.CashierResponseYearTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesByMerchantCacheKey, req.Year, req.MerchantID)
	SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant) ([]*response.CashierResponseMonthSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := GetFromCache[[]*response.CashierResponseMonthSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant, res []*response.CashierResponseMonthSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesByMerchantCacheKey, req.Year, req.MerchantID)
	SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant) ([]*response.CashierResponseYearSales, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := GetFromCache[[]*response.CashierResponseYearSales](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant, res []*response.CashierResponseYearSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesByMerchantCacheKey, req.Year, req.MerchantID)
	SetToCache(ctx, s.store, key, &res, ttlDefault)
}
