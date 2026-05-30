package repo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/domain"
	database "github.com/quangnguyen1505/go-notification-system/internal/notification/infras/postgresql/gen"
	"github.com/quangnguyen1505/go-notification-system/pkg/postgres"
)

var _ domain.NotificationRepo = (*notificationRepo)(nil)

var RepositorySet = wire.NewSet(
	NewNotificationRepo,
	wire.Bind(new(domain.NotificationRepo), new(*notificationRepo)),
)

type notificationRepo struct {
	pg postgres.DBEngine
	q  *database.Queries
}

func NewNotificationRepo(pg postgres.DBEngine) *notificationRepo {
	return &notificationRepo{
		pg: pg,
		q:  database.New(pg.GetDB()),
	}
}

func (r *notificationRepo) Create(ctx context.Context, model *domain.NotificationModel) error {
	if model == nil {
		return errors.New("notification model is nil")
	}

	result, err := r.q.CreateNotification(ctx, database.CreateNotificationParams{
		UserID:  toPgUUID(model.UserID),
		Message: model.Message,
	})
	if err != nil {
		return err
	}
	if result.ID.Valid {
		model.ID = uuid.UUID(result.ID.Bytes)
	}
	return nil
}

func (r *notificationRepo) GetByID(ctx context.Context, model *domain.NotificationModel) (string, error) {
	if model == nil {
		return "", errors.New("notification model is nil")
	}

	result, err := r.q.GetNotificationByID(ctx, toPgUUID(model.ID))
	if err != nil {
		return "", err
	}
	return result.Message, nil
}

func (r *notificationRepo) GetAllByUserID(ctx context.Context, model *domain.NotificationModel) error {
	if model == nil {
		return errors.New("notification model is nil")
	}

	_, err := r.q.GetAllNotificationsByUserID(ctx, toPgUUID(model.UserID))
	return err
}

func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}
