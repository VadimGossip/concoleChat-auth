package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/VadimGossip/concoleChat-auth/internal/api/user"
	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/service"
	serviceMocks "github.com/VadimGossip/concoleChat-auth/internal/service/mocks"
	desc "github.com/VadimGossip/concoleChat-auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = gofakeit.IntRange(1, 2)

		serviceErr = fmt.Errorf("some service error")

		req = &desc.CreateRequest{
			Info: &desc.UserInfo{
				Name:            name,
				Email:           email,
				Password:        password,
				PasswordConfirm: password,
				Role:            desc.Role(role),
			},
		}

		userInfo = &model.UserInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            desc.Role_name[int32(req.Info.Role)],
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
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
				mock.CreateMock.Expect(ctx, userInfo).Return(id, nil)
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
				mock.CreateMock.Expect(ctx, userInfo).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			userServiceMock := test.userServiceMock(mc)

			impl := user.NewImplementation(userServiceMock)
			actualRes, err := impl.Create(test.args.ctx, test.args.req)

			assert.Equal(t, test.want, actualRes)
			assert.Equal(t, test.err, err)
		})
	}
}
