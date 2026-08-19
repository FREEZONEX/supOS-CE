package logicx

import (
	"context"
	"errors"

	"backend/internal/contextx"
	"backend/internal/domain/apikey"
	"backend/internal/domain/asset"
	"backend/internal/domain/audit"
	"backend/internal/domain/flow"
	"backend/internal/domain/iam"
	"backend/internal/domain/uns"
	respx "backend/internal/httpx"
	"backend/internal/repo"
)

func UserID(ctx context.Context) int64 {
	subject, ok := contextx.SubjectFrom(ctx)
	if !ok {
		return 0
	}
	return subject.UserID
}

func Error(err error) error {
	switch {
	case errors.Is(err, uns.ErrInvalid), errors.Is(err, uns.ErrInvalidJSON), errors.Is(err, flow.ErrInvalid),
		errors.Is(err, repo.ErrInvalidArgument), errors.Is(err, repo.ErrFlowFolderCannotDeploy),
		errors.Is(err, asset.ErrInvalid), errors.Is(err, apikey.ErrInvalid),
		errors.Is(err, iam.ErrInvalidRole), errors.Is(err, iam.ErrInvalidUser), errors.Is(err, iam.ErrInvalidPassword):
		return respx.NewHTTPError(400, err.Error())
	case errors.Is(err, repo.ErrSystemReadonly):
		return respx.NewHTTPError(403, err.Error())
	case errors.Is(err, repo.ErrDuplicate), errors.Is(err, repo.ErrUserAccountDuplicate), errors.Is(err, repo.ErrUserEmailDuplicate):
		return respx.NewHTTPError(409, err.Error())
	case errors.Is(err, uns.ErrRestoreConflict):
		return respx.NewHTTPError(409, err.Error())
	case errors.Is(err, repo.ErrNotFound), errors.Is(err, uns.ErrNotFound), errors.Is(err, uns.ErrNotInRecycle),
		errors.Is(err, flow.ErrNotFound), errors.Is(err, asset.ErrNotFound), errors.Is(err, apikey.ErrNotFound),
		errors.Is(err, audit.ErrNotFound):
		return respx.NewHTTPError(404, err.Error())
	default:
		return err
	}
}
