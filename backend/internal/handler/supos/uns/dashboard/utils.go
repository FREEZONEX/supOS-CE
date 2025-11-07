package dashboard

import (
	"backend/internal/common/utils/apiutil"
	"net/http"
)

func getUserID(r *http.Request) string {
	if user := apiutil.GetUserFromContext(r.Context()); user != nil {
		return user.Sub
	}
	return ""
}

func getUsername(r *http.Request) string {
	if user := apiutil.GetUserFromContext(r.Context()); user != nil {
		return user.PreferredUsername
	}
	return ""
}
