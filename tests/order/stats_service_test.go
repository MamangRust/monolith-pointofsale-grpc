package order_test

import (
	"context"
	"testing"
	"time"

	order_cache "github.com/MamangRust/monolith-point-of-sale-order/cache"
	"github.com/MamangRust/monolith-point-of-sale-order/repository"
	"github.com/MamangRust/monolith-point-of-sale-order/service"
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type OrderStatsServiceTestSuite struct {
	tests.BaseTestSuite
	svc           service.OrderStatsService
	svcByMerchant service.OrderStatByMerchantService
	merchantID    int
	userID        int
}

func (s *OrderStatsServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	queries := db.New(s.DBPool())
	repos := repository.NewRepositories(
		queries,
		pb.NewCashierServiceClient(s.Conns["cashier"]),
		pb.NewMerchantServiceClient(s.Conns["merchant"]),
		pb.NewProductServiceClient(s.Conns["product"]),
		pb.NewOrderItemServiceClient(s.Conns["order-item"]),
	)

	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	mencache := order_cache.NewMencache(cacheStore)

	fullSvc := service.NewService(&service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	s.svc = fullSvc.OrderStats
	s.svcByMerchant = fullSvc.OrderStatsByMerchant

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	_, err := s.DBPool().Exec(ctx, "UPDATE orders SET created_at = $1 WHERE order_id = $2",
		time.Now(), orderID)
	s.Require().NoError(err)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyTotalRevenue() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthTotalRevenue{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	res, err := s.svc.FindMonthlyTotalRevenue(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyTotalRevenue() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyTotalRevenue(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyOrder() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindMonthlyOrder(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyOrder() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyOrder(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthTotalRevenueMerchant{
		Year:       now.Year(),
		Month:      int(now.Month()),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindMonthlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyTotalRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearTotalRevenueMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindYearlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyOrderByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthOrderMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindMonthlyOrderByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyOrderByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearOrderMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindYearlyOrderByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func TestOrderStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsServiceTestSuite))
}
