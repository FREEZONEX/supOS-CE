package enums

var DefaultAllowURIs = []string{
	"/logo",
	"/default",
	"/inter-api/supos",
	"/fuxa",
	"/swagger-ui",
	"/assets",
	"/hasura",
	"/chat2db/api",
	"/chat2db/home",
	"/nodered/home",
	"/404",
	"/todo",
	"/plugin",
	"/mf-manifest.json",
	"/403",
	"/copilotkit",
	"/portainer/home",
	"/konga/home",
	"/marimo/home",
	"/grafana/home",
	"/eventflow/home",
	"/emqx/home",
}

func IsDefaultCommonURI(path string) bool {
	for _, uri := range DefaultAllowURIs {
		if path == uri || len(path) > len(uri) && path[:len(uri)+1] == uri+"/" {
			return true
		}
	}
	return false
}
