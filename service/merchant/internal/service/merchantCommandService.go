package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	mencache "github.com/MamangRust/monolith-point-of-sale-merchant/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-merchant/internal/repository"
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/email"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantCommandDeps struct {
	Kafka           *kafka.Kafka
	Cache           mencache.MerchantCommandCache
	UserQuery       repository.UserQueryRepository
	MerchantQuery   repository.MerchantQueryRepository
	MerchantCommand repository.MerchantCommandRepository
	Logger          logger.LoggerInterface
	Observability   observability.TraceLoggerObservability
}

type merchantCommandService struct {
	kafka                     *kafka.Kafka
	mencache                  mencache.MerchantCommandCache
	userRepository            repository.UserQueryRepository
	merchantQueryRepository   repository.MerchantQueryRepository
	merchantCommandRepository repository.MerchantCommandRepository
	logger                    logger.LoggerInterface
	observability             observability.TraceLoggerObservability
}

func NewMerchantCommandService(params *merchantCommandDeps) MerchantCommandService {
	return &merchantCommandService{
		kafka:                     params.Kafka,
		mencache:                  params.Cache,
		userRepository:            params.UserQuery,
		merchantQueryRepository:   params.MerchantQuery,
		merchantCommandRepository: params.MerchantCommand,
		logger:                    params.Logger,
		observability:             params.Observability,
	}
}

func (s *merchantCommandService) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.Merchant, error) {
	const method = "CreateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", request.UserID),
	)
	defer func() {
		end(status)
	}()

	user, err := s.userRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			user_errors.ErrUserNotFoundRes.WithInternal(err),
			method,
			span,
			zap.Int("user.id", request.UserID),
		)
	}

	res, err := s.merchantCommandRepository.CreateMerchant(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedCreateMerchant.WithInternal(err),
			method,
			span,
			zap.Int("user.id", request.UserID),
		)
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Your merchant account has been created successfully. To continue, please upload the required documents for verification. Once completed, our team will review and activate your account.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.UserID),
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": "Initial Verification - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedSendEmail.WithInternal(err),
			method,
			span,
			zap.Int("user.id", int(user.UserID)),
		)
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-created", strconv.Itoa(int(res.MerchantID)), payloadBytes)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedSendEmail.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully created merchant", zap.Int("merchant.id", int(res.MerchantID)))
	return res, nil
}

func (s *merchantCommandService) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.Merchant, error) {
	const method = "UpdateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", *request.MerchantID),
	)
	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.UpdateMerchant(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedUpdateMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchant(ctx, *request.MerchantID)
	logSuccess("Successfully updated merchant", zap.Int("merchant.id", *request.MerchantID))
	return res, nil
}

func (s *merchantCommandService) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error) {
	const method = "UpdateMerchantStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", *request.MerchantID),
	)
	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(ctx, *request.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById.WithInternal(err),
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	user, err := s.userRepository.FindById(ctx, int(merchant.UserID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			user_errors.ErrUserNotFoundRes.WithInternal(err),
			method,
			span,
			zap.Int("user.id", int(merchant.UserID)),
		)
	}

	res, err := s.merchantCommandRepository.UpdateMerchantStatus(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedUpdateMerchant.WithInternal(err),
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	statusReq := request.Status
	subject := ""
	message := ""
	buttonLabel := "Go to Portal"
	link := fmt.Sprintf("https://sanedge.example.com/merchant/%d/dashboard", *request.MerchantID)

	switch statusReq {
	case "active":
		subject = "Your Merchant Account is Now Active"
		message = "Congratulations! Your merchant account has been verified and is now <b>active</b>. You can now fully access all features in the SanEdge Merchant Portal."
	case "inactive":
		subject = "Merchant Account Set to Inactive"
		message = "Your merchant account status has been set to <b>inactive</b>. Please contact support if you believe this is a mistake."
	case "rejected":
		subject = "Merchant Account Rejected"
		message = "We're sorry to inform you that your merchant account has been <b>rejected</b>. Please contact support or review your submissions."
	default:
		return nil, nil
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   subject,
		"Message": message,
		"Button":  buttonLabel,
		"Link":    link,
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": subject,
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedSendEmail.WithInternal(err),
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-update-status", strconv.Itoa(*request.MerchantID), payloadBytes)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedSendEmail.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedMerchant(ctx, *request.MerchantID)
	logSuccess("Successfully updated merchant status", zap.Int("merchant.id", *request.MerchantID))
	return res, nil
}

func (s *merchantCommandService) TrashedMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	const method = "TrashedMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", merchant_id),
	)
	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.TrashedMerchant(ctx, merchant_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedTrashMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully trashed merchant", zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *merchantCommandService) RestoreMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	const method = "RestoreMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", merchant_id),
	)
	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.RestoreMerchant(ctx, merchant_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Merchant](
			s.logger,
			merchant_errors.ErrFailedRestoreMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully restored merchant", zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *merchantCommandService) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	const method = "DeleteMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", merchant_id),
	)
	defer func() {
		end(status)
	}()

	success, err := s.merchantCommandRepository.DeleteMerchantPermanent(ctx, merchant_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedDeleteMerchantPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully deleted merchant permanently", zap.Int("merchant.id", merchant_id), zap.Bool("success", success))
	return success, nil
}

func (s *merchantCommandService) RestoreAllMerchant(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.merchantCommandRepository.RestoreAllMerchant(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedRestoreAllMerchants.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully restored all merchants", zap.Bool("success", success))
	return success, nil
}

func (s *merchantCommandService) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.merchantCommandRepository.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			merchant_errors.ErrFailedDeleteAllMerchantsPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully deleted all merchants permanently", zap.Bool("success", success))
	return success, nil
}
