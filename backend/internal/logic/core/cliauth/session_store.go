package cliauth

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"backend/internal/types"
)

const cliAuthSessionTTL = 10 * time.Minute

type cliAuthSession struct {
	status        string
	apiKey        string
	workspaceName string
	expiresAt     int64
}

var cliAuthSessions sync.Map

func normalizeCliAuthSetupCode(setupCode string) string {
	return strings.ToUpper(strings.TrimSpace(setupCode))
}

func storeCliAuthCompleted(setupCode string, result types.CliAuthBindResp) {
	key := normalizeCliAuthSetupCode(setupCode)
	if key == "" {
		return
	}
	cliAuthSessions.Store(key, cliAuthSession{
		status:        "completed",
		apiKey:        result.ApiKey,
		workspaceName: result.WorkspaceName,
		expiresAt:     time.Now().UTC().Add(cliAuthSessionTTL).UnixMilli(),
	})
}

func loadCliAuthStatus(setupCode string) types.CliAuthStatusResp {
	key := normalizeCliAuthSetupCode(setupCode)
	now := time.Now().UTC()
	if value, ok := cliAuthSessions.Load(key); ok {
		session, _ := value.(cliAuthSession)
		if session.expiresAt > now.UnixMilli() {
			return types.CliAuthStatusResp{
				Status:        session.status,
				ApiKey:        session.apiKey,
				WorkspaceName: session.workspaceName,
				ExpiresAt:     strconv.FormatInt(session.expiresAt, 10),
			}
		}
		cliAuthSessions.Delete(key)
		return types.CliAuthStatusResp{
			Status:        "expired",
			WorkspaceName: "default",
			ExpiresAt:     strconv.FormatInt(session.expiresAt, 10),
		}
	}
	return types.CliAuthStatusResp{
		Status:        "pending",
		WorkspaceName: "default",
		ExpiresAt:     strconv.FormatInt(now.Add(cliAuthSessionTTL).UnixMilli(), 10),
	}
}
