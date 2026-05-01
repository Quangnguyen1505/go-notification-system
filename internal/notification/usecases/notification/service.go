package notification

import (
	"context"

	"github.com/google/wire"
	"github.com/quangnguyen1505/go-notification-system/internal/notification/domain"
)

var _ UseCase = (*service)(nil)

var UseCaseSet = wire.NewSet(
	NewService,
	wire.Bind(new(UseCase), new(*service)),
)

type service struct {
	repo domain.NotificationRepo
}

func NewService(repo domain.NotificationRepo) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, model *domain.NotificationModel) error {
	return s.repo.Create(ctx, model)
}

func (s *service) GetByID(ctx context.Context, model *domain.NotificationModel) error {
	return s.repo.GetByID(ctx, model)
}

func (s *service) GetAllByUserID(ctx context.Context, model *domain.NotificationModel) error {
	return s.repo.GetAllByUserID(ctx, model)
}
