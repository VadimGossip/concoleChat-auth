package tests

//не могу понять как мокать TxManager ругается на expected f

//import (
//	"context"
//	"fmt"
//	"github.com/VadimGossip/concoleChat-auth/internal/model"
//	"github.com/VadimGossip/concoleChat-auth/internal/repository"
//	repoMocks "github.com/VadimGossip/concoleChat-auth/internal/repository/mocks"
//	"github.com/VadimGossip/concoleChat-auth/internal/service"
//	servMocks "github.com/VadimGossip/concoleChat-auth/internal/service/mocks"
//	userService "github.com/VadimGossip/concoleChat-auth/internal/service/user"
//	db "github.com/VadimGossip/platform_common/pkg/db/postgres"
//	"github.com/brianvoe/gofakeit/v6"
//	"github.com/gojuno/minimock/v3"
//	"github.com/stretchr/testify/assert"
//	"testing"
//	"time"
//)
//
//func TestCreate(t *testing.T) {
//	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
//	type txManagerMockFunc func(f any, mc *minimock.Controller) db.TxManager
//	type auditServiceMockFunc func(mc *minimock.Controller) service.AuditService
//	type userCacheServiceMockFunc func(mc *minimock.Controller) service.UserCacheService
//
//	type args struct {
//		ctx context.Context
//		req *model.UserInfo
//	}
//
//	var (
//		ctx = context.Background()
//		//repoErr = fmt.Errorf("repo error")
//
//		id = gofakeit.Int64()
//
//		req = &model.UserInfo{
//			Name:            gofakeit.Name(),
//			Email:           gofakeit.Email(),
//			Password:        "pwd",
//			PasswordConfirm: "pwd",
//			Role:            model.AdminRole,
//		}
//
//		user = &model.User{
//			ID:        id,
//			Info:      *req,
//			CreatedAt: time.Now(),
//			UpdatedAt: nil,
//		}
//
//		auditMsg = &model.Audit{
//			Action:     "create user",
//			CallParams: fmt.Sprintf("info %+v", req),
//		}
//	)
//
//	tests := []struct {
//		caseName             string
//		args                 args
//		expectedResult       int64
//		expectedError        error
//		userRepositoryMock   userRepositoryMockFunc
//		txManagerMock        txManagerMockFunc
//		userCacheServiceMock userCacheServiceMockFunc
//		auditServiceMock     auditServiceMockFunc
//	}{
//		{
//			caseName: "success case",
//			args: args{
//				ctx: ctx,
//				req: req,
//			},
//			expectedResult: id,
//			expectedError:  nil,
//			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
//				mock := repoMocks.NewUserRepositoryMock(mc)
//				mock.CreateMock.Expect(ctx, req).Return(user, nil)
//				return mock
//			},
//			txManagerMock: func(f any, mc *minimock.Controller) db.TxManager {
//				mock := repoMocks.NewTxManagerMock(mc)
//				mock.ReadCommittedMock.Expect(ctx, f).Return(nil)
//				return mock
//			},
//			userCacheServiceMock: func(mc *minimock.Controller) service.UserCacheService {
//				mock := servMocks.NewUserCacheServiceMock(mc)
//				mock.SetMock.Expect(ctx, user).Times(1).Return(nil)
//				return mock
//			},
//			auditServiceMock: func(mc *minimock.Controller) service.AuditService {
//				mock := servMocks.NewAuditServiceMock(mc)
//				mock.CreateMock.Expect(ctx, auditMsg).Times(1).Return(nil)
//				return mock
//			},
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.caseName, func(t *testing.T) {
//			mc := minimock.NewController(t)
//			userRepositoryMock := test.userRepositoryMock(mc)
//			auditServiceMock := test.auditServiceMock(mc)
//			userCacheServiceMock := test.userCacheServiceMock(mc)
//			txManagerMock := test.txManagerMock(func(ctx context.Context) error {
//				var txErr error
//				user, txErr = userRepositoryMock.Create(ctx, req)
//				if txErr != nil {
//					return txErr
//				}
//
//				txErr = auditServiceMock.Create(ctx, auditMsg)
//				if txErr != nil {
//					return txErr
//				}
//
//				return nil
//			}, mc)
//
//			s := userService.NewService(userRepositoryMock,
//				userCacheServiceMock,
//				auditServiceMock,
//				txManagerMock)
//			actualId, err := s.Create(test.args.ctx, test.args.req)
//
//			assert.Equal(t, test.expectedResult, actualId)
//			assert.Equal(t, test.expectedError, err)
//		})
//	}
//}
