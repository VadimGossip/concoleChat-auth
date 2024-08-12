package user

import (
	"context"

	"github.com/VadimGossip/concoleChat-auth/internal/converter"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if err := i.userService.Update(ctx, req.Id, converter.ToUpdateUserInfoFromDesc(req.GetInfo())); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
