package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	repoMocks "github.com/VadimGossip/concoleChat-auth/internal/repository/mocks"
	"github.com/VadimGossip/concoleChat-auth/internal/service/audit"
)

func TestCreate(t *testing.T) {
	type auditRepositoryMockFunc func(mc *minimock.Controller) repository.AuditRepository

	type args struct {
		ctx context.Context
		req *model.Audit
	}

	var (
		ctx     = context.Background()
		repoErr = fmt.Errorf("repo error")

		req = &model.Audit{
			Action:     "some action",
			CallParams: "some action params",
		}
	)

	tests := []struct {
		caseName            string
		args                args
		expectedError       error
		auditRepositoryMock auditRepositoryMockFunc
	}{
		{
			caseName: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			expectedError: nil,
			auditRepositoryMock: func(mc *minimock.Controller) repository.AuditRepository {
				mock := repoMocks.NewAuditRepositoryMock(mc)
				mock.CreateMock.Times(1).Expect(ctx, req).Return(nil)
				return mock
			},
		},
		{
			caseName: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			expectedError: repoErr,
			auditRepositoryMock: func(mc *minimock.Controller) repository.AuditRepository {
				mock := repoMocks.NewAuditRepositoryMock(mc)
				mock.CreateMock.Times(1).Expect(ctx, req).Return(repoErr)
				return mock
			},
		},
	}
	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			mc := minimock.NewController(t)
			mockAuditRepository := test.auditRepositoryMock(mc)

			s := audit.NewService(mockAuditRepository)
			err := s.Create(ctx, test.args.req)

			assert.Equal(t, test.expectedError, err)
		})
	}
}
