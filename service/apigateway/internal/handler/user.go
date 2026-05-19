package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	gateway_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/internal/redis/api/gateway_cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandleApi struct {
	client     pb.UserServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.UserResponseMapper
	apiHandler errors.ApiHandler
	cache      *gateway_cache.GatewayCache
}

func NewHandlerUser(
	router *echo.Echo,
	client pb.UserServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.UserResponseMapper,
	apiHandler errors.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *userHandleApi {
	userHandler := &userHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerUser := router.Group("/api/user")

	routerUser.GET("", apiHandler.Handle("get-user-findalluser", userHandler.FindAllUser))
	routerUser.GET("/:id", apiHandler.Handle("get-user-findbyid", userHandler.FindById))
	routerUser.GET("/active", apiHandler.Handle("get-user-findbyactive", userHandler.FindByActive))
	routerUser.GET("/trashed", apiHandler.Handle("get-user-findbytrashed", userHandler.FindByTrashed))

	routerUser.POST("/create", apiHandler.Handle("post-user-create", userHandler.Create))
	routerUser.POST("/update/:id", apiHandler.Handle("post-user-update", userHandler.Update))

	routerUser.POST("/trashed/:id", apiHandler.Handle("post-user-trasheduser", userHandler.TrashedUser))
	routerUser.POST("/restore/:id", apiHandler.Handle("post-user-restoreuser", userHandler.RestoreUser))
	routerUser.DELETE("/permanent/:id", apiHandler.Handle("delete-user-deleteuserpermanent", userHandler.DeleteUserPermanent))

	routerUser.POST("/restore/all", apiHandler.Handle("post-user-restorealluser", userHandler.RestoreAllUser))
	routerUser.POST("/permanent/all", apiHandler.Handle("delete-user-deletealluserpermanent", userHandler.DeleteAllUserPermanent))

	return userHandler
}

func (h *userHandleApi) FindAllUser(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("user:findalluser:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationUser](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationUser(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid user ID")
	}

	cacheKey := fmt.Sprintf("user:findbyid:id_%d", id)
	if cached, found := gateway_cache.Get[response.ApiResponseUser](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseUser(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) FindByActive(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("user:findbyactive:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationUserDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationUserDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) FindByTrashed(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("user:findbytrashed:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationUserDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationUserDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var body requests.CreateUserRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.CreateUserRequest{
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseUser(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "user:*")
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid user ID")
	}

	var body requests.UpdateUserRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.UpdateUserRequest{
		Id:              int32(id),
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseUser(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "user:*")
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) TrashedUser(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid user ID")
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedUser(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseUserDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "user:*")
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) RestoreUser(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid user ID")
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreUser(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseUserDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "user:*")
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) DeleteUserPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid user ID")
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteUserPermanent(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseUserDelete(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "user:*")
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) RestoreAllUser(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllUser(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseUserAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "user:*")
	return c.JSON(http.StatusOK, so)
}

func (h *userHandleApi) DeleteAllUserPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllUserPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseUserAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "user:*")
	return c.JSON(http.StatusOK, so)
}
