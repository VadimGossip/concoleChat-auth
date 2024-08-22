package consumer

import (
	"context"
)

type UserConsumerService interface {
	RunConsumer(ctx context.Context) error
}
