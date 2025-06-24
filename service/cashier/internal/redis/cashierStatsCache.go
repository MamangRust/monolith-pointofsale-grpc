package mencache

import (
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

const (
	cashierStatsMonthTotalSalesCacheKey = "cashier:stats:month:%d:year:%d"
	cashierStatsYearTotalSalesCacheKey  = "cashier:stats:year:%d"

	cashierStatsMonthSalesCacheKey = "cashier:stats:month:%d"
	cashierStatsYearSalesCacheKey  = "cashier:stats:year:%d"
)

type cashierStatsCache struct {
	store *CacheStore
}

func NewCashierStatsCache(store *CacheStore) *cashierStatsCache {
	return &cashierStatsCache{store: store}
}

func (s *cashierStatsCache) GetMonthlyTotalSalesCache(req *requests.MonthTotalSales) ([]*response.CashierResponseMonthTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesCacheKey, req.Month, req.Year)
	result, found := GetFromCache[[]*response.CashierResponseMonthTotalSales](s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetMonthlyTotalSalesCache(req *requests.MonthTotalSales, res []*response.CashierResponseMonthTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesCacheKey, req.Month, req.Year)
	SetToCache(s.store, key, &res, ttlDefault)
}

func (s *cashierStatsCache) GetYearlyTotalSalesCache(year int) ([]*response.CashierResponseYearTotalSales, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesCacheKey, year)
	result, found := GetFromCache[[]*response.CashierResponseYearTotalSales](s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetYearlyTotalSalesCache(year int, res []*response.CashierResponseYearTotalSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesCacheKey, year)
	SetToCache(s.store, key, &res, ttlDefault)
}

// Get & Set MonthlySales
func (s *cashierStatsCache) GetMonthlySalesCache(year int) ([]*response.CashierResponseMonthSales, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesCacheKey, year)
	result, found := GetFromCache[[]*response.CashierResponseMonthSales](s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetMonthlySalesCache(year int, res []*response.CashierResponseMonthSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesCacheKey, year)
	SetToCache(s.store, key, &res, ttlDefault)
}

func (s *cashierStatsCache) GetYearlySalesCache(year int) ([]*response.CashierResponseYearSales, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesCacheKey, year)
	result, found := GetFromCache[[]*response.CashierResponseYearSales](s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetYearlySalesCache(year int, res []*response.CashierResponseYearSales) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesCacheKey, year)
	SetToCache(s.store, key, &res, ttlDefault)
}
