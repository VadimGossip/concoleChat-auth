package user

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/converter"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}
