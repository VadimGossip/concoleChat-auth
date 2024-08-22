package producer

import (
	"github.com/VadimGossip/concoleChat-auth/internal/model"
)

type UserProducerService interface {
	ProduceCreate(info *model.UserInfo) error
}
