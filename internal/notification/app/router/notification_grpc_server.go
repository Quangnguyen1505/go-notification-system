package router

import (
	"context"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/domain"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/usecases/notification"
	"github.com/quangnguyen1505/go-notification-system/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var _ gen.NotificationServiceServer = (*notificationGRPCServer)(nil)

var NotificationGRPCServerSet = wire.NewSet(NewNotificationGRPCServer)

type notificationGRPCServer struct {
	gen.UnimplementedNotificationServiceServer
	uc notification.UseCase
}

func NewNotificationGRPCServer(
	grpcServer *grpc.Server,
	uc notification.UseCase,
) gen.NotificationServiceServer {
	svc := notificationGRPCServer{
		uc: uc,
	}

	gen.RegisterNotificationServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (s *notificationGRPCServer) CreateNotification(
	ctx context.Context,
	req *gen.CreateNotificationRequest,
) (*gen.CreateNotificationResponse, error) {
	model := &domain.NotificationModel{
		ID: uuid.New(),
	}

	if req.GetUserId() != "" {
		userID, err := uuid.Parse(req.GetUserId())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid user_id")
		}
		model.UserID = userID
	}

	model.Message = req.GetSubject() + ": " + req.GetBody()

	if err := s.uc.Create(ctx, model); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.CreateNotificationResponse{
		NotificationId: model.ID.String(),
		Status:         "created",
	}, nil
}

func (s *notificationGRPCServer) GetNotification(
	ctx context.Context,
	req *gen.GetNotificationRequest,
) (*gen.GetNotificationResponse, error) {
	notificationID, err := uuid.Parse(req.GetNotificationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid notification_id")
	}

	model := &domain.NotificationModel{ID: notificationID}
	noti_id, err := s.uc.GetByID(ctx, model)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetNotificationResponse{
		NotificationId: noti_id,
		Status:         "ok",
	}, nil
}

func (s *notificationGRPCServer) BatchNotification(
	ctx context.Context,
	req *gen.BatchNotificationRequest,
) (*gen.BatchNotificationResponse, error) {
	results := make([]*gen.BatchNotificationResult, 0, len(req.GetNotifications()))

	for _, n := range req.GetNotifications() {
		model := &domain.NotificationModel{
			ID:      uuid.New(),
			Message: n.GetSubject() + ": " + n.GetBody(),
		}

		if n.GetUserId() != "" {
			userID, err := uuid.Parse(n.GetUserId())
			if err != nil {
				results = append(results, &gen.BatchNotificationResult{
					NotificationId: "",
					Status:         "error",
					Error:          "invalid user_id",
				})
				continue
			}
			model.UserID = userID
		}

		if err := s.uc.Create(ctx, model); err != nil {
			results = append(results, &gen.BatchNotificationResult{
				NotificationId: model.ID.String(),
				Status:         "error",
				Error:          err.Error(),
			})
			continue
		}

		results = append(results, &gen.BatchNotificationResult{
			NotificationId: model.ID.String(),
			Status:         "created",
			Error:          "",
		})
	}

	return &gen.BatchNotificationResponse{Results: results}, nil
}
