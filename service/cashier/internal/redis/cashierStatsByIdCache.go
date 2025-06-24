package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

const (
	cashierStatsMonthTotalSalesByIdCacheKey = "cashier:stats:month:%d:year:%d:id:%d"
	cashierStatsYearTotalSalesByIdCacheKey  = "cashier:stats:year:%d:id:%d"

	cashierStatsMonthSalesByIdCacheKey = "cashier:stats:month:%d:id:%d"
	cashierStatsYearSalesByIdCacheKey  = "cashier:stats:year:%d:id:%d"
)

type cashierStatsByIdCache struct {
	store *CacheStore
}

func NewCashierStatsByIdCache(store *CacheStore) *cashierStatsByIdCache {
	return &cashierStatsByIdCache{store: store}
}

func (s *cashierStatsByIdCache) GetMonthlyTotalSalesByIdCache(req *requests.MonthTotalSalesCashier) ([]*response.CashierResponseMonthTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByIdCacheKey, req.Month, req.Year, req.CashierID)
	result, found := GetFromCache[[]*response.CashierResponseMonthTotalSales](s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetMonthlyTotalSalesByIdCache(req *requests.MonthTotalSalesCashier, res []*response.CashierResponseMonthTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByIdCacheKey, req.Month, req.Year, req.CashierID)
	SetToCache(s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetYearlyTotalSalesByIdCache(req *requests.YearTotalSalesCashier) ([]*response.CashierResponseYearTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := GetFromCache[[]*response.CashierResponseYearTotalSales](s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetYearlyTotalSalesByIdCache(req *requests.YearTotalSalesCashier, res []*response.CashierResponseYearTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesByIdCacheKey, req.Year, req.CashierID)
	SetToCache(s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetMonthlyCashierByIdCache(req *requests.MonthCashierId) ([]*response.CashierResponseMonthSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := GetFromCache[[]*response.CashierResponseMonthSales](s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetMonthlyCashierByIdCache(req *requests.MonthCashierId, res []*response.CashierResponseMonthSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesByIdCacheKey, req.Year, req.CashierID)
	SetToCache(s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetYearlyCashierByIdCache(req *requests.YearCashierId) ([]*response.CashierResponseYearSales, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := GetFromCache[[]*response.CashierResponseYearSales](s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetYearlyCashierByIdCache(req *requests.YearCashierId, res []*response.CashierResponseYearSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesByIdCacheKey, req.Year, req.CashierID)
	SetToCache(s.store, key, &res, ttlDefault)
}
