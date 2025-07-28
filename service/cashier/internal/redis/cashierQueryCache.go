package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

const (
	cashierAllCacheKey     = "cashier:all:page:%d:pageSize:%d:search:%s"
	cashierByIdCacheKey    = "cashier:id:%d"
	cashierActiveCacheKey  = "cashier:active:page:%d:pageSize:%d:search:%s"
	cashierTrashedCacheKey = "cashier:trashed:page:%d:pageSize:%d:search:%s"

	cashierByMerchantCacheKey = "cashier:merchant:%d:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type cashierCacheResponse struct {
	Data         []*response.CashierResponse `json:"data"`
	TotalRecords *int                        `json:"totalRecords"`
}

type cashierCacheResponseDeleteAt struct {
	Data         []*response.CashierResponseDeleteAt `json:"data"`
	TotalRecords *int                                `json:"totalRecords"`
}

type cashierQueryCache struct {
	store *CacheStore
}

func NewCashierQueryCache(store *CacheStore) *cashierQueryCache {
	return &cashierQueryCache{store: store}
}

func (s *cashierQueryCache) GetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponse, *int, bool) {
	key := fmt.Sprintf(cashierAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[cashierCacheResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *cashierQueryCache) SetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers, data []*response.CashierResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.CashierResponse{}
	}

	key := fmt.Sprintf(cashierAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &cashierCacheResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*response.CashierResponse, *int, bool) {
	key := fmt.Sprintf(cashierByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[cashierCacheResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *cashierQueryCache) SetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant, data []*response.CashierResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.CashierResponse{}
	}

	key := fmt.Sprintf(cashierByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &cashierCacheResponse{Data: data, TotalRecords: total}
	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(cashierActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[cashierCacheResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *cashierQueryCache) SetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers, data []*response.CashierResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CashierResponseDeleteAt{}
	}

	key := fmt.Sprintf(cashierActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &cashierCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(cashierTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := GetFromCache[cashierCacheResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *cashierQueryCache) SetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers, data []*response.CashierResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.CashierResponseDeleteAt{}
	}

	key := fmt.Sprintf(cashierTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &cashierCacheResponseDeleteAt{Data: data, TotalRecords: total}
	SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *cashierQueryCache) GetCachedCashier(ctx context.Context, id int) (*response.CashierResponse, bool) {
	key := fmt.Sprintf(cashierByIdCacheKey, id)
	result, found := GetFromCache[*response.CashierResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *cashierQueryCache) SetCachedCashier(ctx context.Context, data *response.CashierResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(cashierByIdCacheKey, data.ID)

	SetToCache(ctx, s.store, key, data, ttlDefault)
}
