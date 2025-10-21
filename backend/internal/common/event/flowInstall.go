package event

// Constants for flow operations
const (
	FlowOperationInstall   = "INSTALL"
	FlowOperationUninstall = "UNINSTALL"
)

// FlowInstallEvent defines an event for flow installation or uninstallation.
type FlowInstallEvent struct {
	FlowName  string
	Operation string
}
