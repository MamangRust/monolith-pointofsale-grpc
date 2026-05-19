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

type roleHandleApi struct {
	client     pb.RoleServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.RoleResponseMapper
	apiHandler errors.ApiHandler
	cache      *gateway_cache.GatewayCache
}

func NewHandlerRole(
	router *echo.Echo,
	client pb.RoleServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.RoleResponseMapper,
	apiHandler errors.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *roleHandleApi {
	roleHandler := &roleHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerRole := router.Group("/api/role")

	routerRole.GET("", apiHandler.Handle("get-role-findall", roleHandler.FindAll))
	routerRole.GET("/:id", apiHandler.Handle("get-role-findbyid", roleHandler.FindById))
	routerRole.GET("/active", apiHandler.Handle("get-role-findbyactive", roleHandler.FindByActive))
	routerRole.GET("/trashed", apiHandler.Handle("get-role-findbytrashed", roleHandler.FindByTrashed))
	routerRole.GET("/user/:user_id", apiHandler.Handle("get-role-findbyuserid", roleHandler.FindByUserId))
	routerRole.POST("", apiHandler.Handle("post-role-create", roleHandler.Create))
	routerRole.POST("/update/:id", apiHandler.Handle("post-role-update", roleHandler.Update))
	routerRole.POST("/trashed/:id", apiHandler.Handle("post-role-trashed", roleHandler.Trashed))
	routerRole.POST("/restore/:id", apiHandler.Handle("post-role-restore", roleHandler.Restore))
	routerRole.DELETE("/permanent/:id", apiHandler.Handle("delete-role-deletepermanent", roleHandler.DeletePermanent))
	routerRole.POST("/restore/all", apiHandler.Handle("post-role-restoreall", roleHandler.RestoreAll))
	routerRole.DELETE("/permanent-all", apiHandler.Handle("delete-role-deleteallpermanent", roleHandler.DeleteAllPermanent))

	return roleHandler
}

func (h *roleHandleApi) FindAll(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("role:findall:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationRole](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllRole(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationRole(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return errors.NewBadRequestError("invalid role ID")
	}

	cacheKey := fmt.Sprintf("role:findbyid:id_%d", roleID)
	if cached, found := gateway_cache.Get[response.ApiResponseRole](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.client.FindByIdRole(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseRole(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) FindByActive(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("role:findbyactive:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationRoleDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationRoleDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) FindByTrashed(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("role:findbytrashed:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationRoleDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationRoleDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) FindByUserId(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		return errors.NewBadRequestError("invalid user ID")
	}

	cacheKey := fmt.Sprintf("role:findbyuserid:user_id_%d", userID)
	if cached, found := gateway_cache.Get[response.ApiResponsesRole](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindByIdUserRoleRequest{
		UserId: int32(userID),
	}

	res, err := h.client.FindByUserId(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsesRole(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req requests.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := req.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	reqPb := &pb.CreateRoleRequest{
		Name: req.Name,
	}

	res, err := h.client.CreateRole(ctx, reqPb)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseRole(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "role:*")
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) Update(c echo.Context) error {
	ctx := c.Request().Context()
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return errors.NewBadRequestError("invalid role ID")
	}

	var req requests.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := req.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	reqPb := &pb.UpdateRoleRequest{
		Id:   int32(roleID),
		Name: req.Name,
	}

	res, err := h.client.UpdateRole(ctx, reqPb)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseRole(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "role:*")
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) Trashed(c echo.Context) error {
	ctx := c.Request().Context()
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return errors.NewBadRequestError("invalid role ID")
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.client.TrashedRole(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseRole(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "role:*")
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) Restore(c echo.Context) error {
	ctx := c.Request().Context()
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return errors.NewBadRequestError("invalid role ID")
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.client.RestoreRole(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseRole(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "role:*")
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) DeletePermanent(c echo.Context) error {
	ctx := c.Request().Context()
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return errors.NewBadRequestError("invalid role ID")
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.client.DeleteRolePermanent(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseRoleDelete(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "role:*")
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllRole(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseRoleAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "role:*")
	return c.JSON(http.StatusOK, so)
}

func (h *roleHandleApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseRoleAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "role:*")
	return c.JSON(http.StatusOK, so)
}
