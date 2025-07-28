package handler

import (
	"fmt"
	"strconv"

	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-pkg/upload_image"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
)

type ServiceConnections struct {
	Auth        *grpc.ClientConn
	Role        *grpc.ClientConn
	User        *grpc.ClientConn
	Cashier     *grpc.ClientConn
	Category    *grpc.ClientConn
	Merchant    *grpc.ClientConn
	OrderItem   *grpc.ClientConn
	Order       *grpc.ClientConn
	Product     *grpc.ClientConn
	Transaction *grpc.ClientConn
}

type Deps struct {
	Token              auth.TokenManager
	E                  *echo.Echo
	Logger             logger.LoggerInterface
	Mapping            *response_api.ResponseApiMapper
	ImageUpload        upload_image.ImageUploads
	ServiceConnections *ServiceConnections
}

func NewHandler(deps *Deps) {
	clientAuth := pb.NewAuthServiceClient(deps.ServiceConnections.Auth)
	clientRole := pb.NewRoleServiceClient(deps.ServiceConnections.Role)
	clientUser := pb.NewUserServiceClient(deps.ServiceConnections.User)
	clientCategory := pb.NewCategoryServiceClient(deps.ServiceConnections.Category)
	clientCashier := pb.NewCashierServiceClient(deps.ServiceConnections.Cashier)
	clientMerchant := pb.NewMerchantServiceClient(deps.ServiceConnections.Merchant)
	clientMerchantDocument := pb.NewMerchantDocumentServiceClient(deps.ServiceConnections.Merchant)
	clientOrderItem := pb.NewOrderItemServiceClient(deps.ServiceConnections.OrderItem)
	clientOrder := pb.NewOrderServiceClient(deps.ServiceConnections.Order)
	clientProduct := pb.NewProductServiceClient(deps.ServiceConnections.Product)
	clientTransaction := pb.NewTransactionServiceClient(deps.ServiceConnections.Transaction)

	NewHandlerAuth(deps.E, clientAuth, deps.Logger, deps.Mapping.AuthResponseMapper)
	NewHandlerRole(deps.E, clientRole, deps.Logger, deps.Mapping.RoleResponseMapper)
	NewHandlerUser(deps.E, clientUser, deps.Logger, deps.Mapping.UserResponseMapper)
	NewHandlerCategory(deps.E, clientCategory, deps.Logger, deps.Mapping.CategoryResponseMapper)
	NewHandlerCashier(deps.E, clientCashier, deps.Logger, deps.Mapping.CashierResponseMapper)
	NewHandlerMerchant(deps.E, clientMerchant, deps.Logger, deps.Mapping.MerchantResponseMapper)
	NewHandlerMerchantDocument(deps.E, clientMerchantDocument, deps.Logger, deps.Mapping.MerchantDocumentProMapper)
	NewHandlerOrderItem(deps.E, clientOrderItem, deps.Logger, deps.Mapping.OrderItemResponseMapper)
	NewHandlerOrder(deps.E, clientOrder, deps.Logger, deps.Mapping.OrderResponseMapper)
	NewHandlerProduct(deps.E, clientProduct, deps.Logger, deps.Mapping.ProductResponseMapper, deps.ImageUpload)
	NewHandlerTransaction(deps.E, clientTransaction, deps.Logger, deps.Mapping.TransactionResponseMapper)
}

func parseQueryInt(c echo.Context, key string, defaultValue int) int {
	val, err := strconv.Atoi(c.QueryParam(key))
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}

func parseQueryIntWithValidation(c echo.Context, key string, min, max int) (int, error) {
	valStr := c.QueryParam(key)
	val, err := strconv.Atoi(valStr)
	if err != nil || val < min || val > max {
		return 0, fmt.Errorf("invalid %s: %s", key, valStr)
	}
	return val, nil
}

func parseQueryStringRequired(c echo.Context, key string) (string, error) {
	val := c.QueryParam(key)
	if val == "" {
		return "", fmt.Errorf("missing or empty query parameter: %s", key)
	}
	return val, nil
}
