package event

// PluginPreUninstallEvent defines an event triggered before a plugin is uninstalled.
// It contains a callback function to be executed.
type PluginPreUninstallEvent struct {
	PluginName        string
	UninstallCallback func(string)
}
