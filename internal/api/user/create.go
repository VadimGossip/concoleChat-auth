package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/VadimGossip/concoleChat-auth/internal/converter"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (i *Implementation) CreateAsync(_ context.Context, req *desc.CreateAsyncRequest) (*emptypb.Empty, error) {
	err := i.userAsyncService.ProduceCreate(converter.ToUserInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
