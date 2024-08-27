package converter

import (
	"github.com/VadimGossip/concoleChat-auth/internal/model"

	repoModel "github.com/VadimGossip/concoleChat-auth/internal/repository/audit/model"
)

func ToRepoFromAudit(audit *model.Audit) *repoModel.Audit {
	return &repoModel.Audit{
		ID:         audit.ID,
		Action:     audit.Action,
		CallParams: audit.CallParams,
	}
}
