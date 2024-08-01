package audit

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

func (s *service) Create(ctx context.Context, audit *model.Audit) error {
	return s.auditRepository.Create(ctx, audit)
}
