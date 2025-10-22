package event

// Constants for flow operations
const (
	FlowOperationInstall   = "INSTALL"
	FlowOperationUninstall = "UNINSTALL"
)

// FlowInstallEvent defines an event for flow installation or uninstallation.
type FlowInstallEvent struct {
	ApplicationEvent
	FlowName  string
	Operation string
}
