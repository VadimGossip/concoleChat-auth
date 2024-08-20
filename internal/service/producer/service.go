package producer

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

type UserProducerService interface {
	ProduceCreate(ctx context.Context, info *model.UserInfo) error
}
