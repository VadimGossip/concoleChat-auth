package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	serviceMocks "github.com/VadimGossip/concoleChat-auth/internal/service/mocks"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService
	type userAsyncServiceMockFunc func(mc *minimock.Controller) service.UserProducerService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = gofakeit.IntRange(1, 2)

		serviceErr = fmt.Errorf("some service error")

		descRole = desc.Role(role)
		req      = &desc.UpdateRequest{
			Id: id,
			Info: &desc.UpdateUserInfo{
				Name:  wrapperspb.String(name),
				Email: wrapperspb.String(email),
				Role:  &descRole,
			},
		}

		userInfo = &model.UpdateUserInfo{
			Name:  &name,
			Email: &email,
			Role:  desc.Role_name[int32(*req.Info.Role)],
		}

		res = &emptypb.Empty{}
	)

	tests := []struct {
		name                 string
		args                 args
		want                 *emptypb.Empty
		err                  error
		userServiceMock      userServiceMockFunc
		userAsyncServiceMock userAsyncServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, userInfo).Return(nil)
				return mock
			},
			userAsyncServiceMock: func(mc *minimock.Controller) service.UserProducerService {
				mock := serviceMocks.NewUserProducerServiceMock(mc)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, userInfo).Return(serviceErr)
				return mock
			},
			userAsyncServiceMock: func(mc *minimock.Controller) service.UserProducerService {
				mock := serviceMocks.NewUserProducerServiceMock(mc)
				return mock
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			userServiceMock := test.userServiceMock(mc)
			userAsyncServiceMock := test.userAsyncServiceMock(mc)

			impl := user.NewImplementation(userServiceMock, userAsyncServiceMock)
			actualRes, err := impl.Update(test.args.ctx, test.args.req)

			assert.Equal(t, test.want, actualRes)
			assert.Equal(t, test.err, err)
		})
	}
}
