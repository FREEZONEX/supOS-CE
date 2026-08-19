package cachekeys

import "fmt"

const prefix = "edge"

func PermissionByUser(userID int64) string {
	return fmt.Sprintf("%s:permission:user:%d", prefix, userID)
}

func APIKeyByHash(hash string) string {
	return fmt.Sprintf("%s:apikey:%s", prefix, hash)
}

func UNSDefinitionByID(id int64) string {
	return fmt.Sprintf("%s:uns:definition:%d", prefix, id)
}

func GatewayRoutesVersion() string {
	return prefix + ":gateway:routes:version"
}

func GatewayRoutesSnapshot() string {
	return prefix + ":gateway:routes:snapshot"
}

func OutboxLock(eventID string) string {
	return fmt.Sprintf("%s:outbox:lock:%s", prefix, eventID)
}
