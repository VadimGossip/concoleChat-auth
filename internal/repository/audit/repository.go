package audit

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	db "github.com/VadimGossip/platform_common/pkg/db/postgres"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	def "github.com/VadimGossip/concoleChat-auth/internal/repository"
	"github.com/VadimGossip/concoleChat-auth/internal/repository/audit/converter"
)

const (
	auditTableName   string = "audit_log"
	actionColumn     string = "action"
	callParamsColumn string = "call_params"
	repoName         string = "audit_repository"
)

var _ def.AuditRepository = (*repository)(nil)

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, audit *model.Audit) error {
	repoAudit := converter.ToRepoFromAudit(audit)
	userInsert := sq.Insert(auditTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(actionColumn, callParamsColumn).
		Values(repoAudit.Action, repoAudit.CallParams)

	query, args, err := userInsert.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     repoName + ".Create",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
