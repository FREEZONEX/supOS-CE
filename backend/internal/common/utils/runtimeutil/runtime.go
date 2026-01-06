package runtimeutil

import "os"

// getEnv retrieves the value of the "ENV" or "GIN_MODE" environment variable.
// It's a common practice in Go to use environment variables to determine the runtime environment.
func getEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = os.Getenv("GIN_MODE")
	}
	if env == "" {
		// Default to "development" if not set, which is a safe default.
		return "development"
	}
	return env
}

// IsLocalEnv checks if the current environment is "local", "dev", or "development".
func IsLocalEnv() bool {
	env := getEnv()
	switch env {
	case "local", "dev", "development":
		return true
	default:
		return false
	}
}
