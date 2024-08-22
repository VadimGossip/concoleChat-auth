package audit

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	db "github.com/VadimGossip/platform_common/pkg/db/postgres"

	def "github.com/VadimGossip/concoleChat-auth/internal/repository"
)

const (
	accessibleRolesTableName string = "accessible_roles"
	endpointAddressColumn    string = "endpoint_address"
	roleColumn               string = "role"
	repoName                 string = "access_repository"
)

var _ def.AccessRepository = (*repository)(nil)

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AccessibleByRole(ctx context.Context, role, endpointAddress string) (bool, error) {
	countSelect := sq.Select("count *").
		From(accessibleRolesTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{endpointAddressColumn: endpointAddressColumn}).
		Where(sq.Eq{roleColumn: roleColumn})

	query, args, err := countSelect.ToSql()
	if err != nil {
		return false, err
	}

	var count int
	q := db.Query{
		Name:     repoName + ".AccessibleByRole",
		QueryRaw: query,
	}
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}
