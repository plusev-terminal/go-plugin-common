package meta

type Meta struct {
	PluginID    string          `json:"pluginId"`
	Name        string          `json:"name"`
	AppID       string          `json:"appId"`
	Category    string          `json:"category"` // Plugin category specific to the app id.
	Description string          `json:"description"`
	Author      string          `json:"author"`
	Version     string          `json:"version"`
	Repository  string          `json:"repository"`
	Tags        []string        `json:"tags"`
	Contacts    []AuthorContact `json:"contacts"`
	Resources   ResourceAccess  `json:"resources"`
	Features    []string        `json:"features"` // List of supported features
}

type AuthorContact struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type ResourceAccess struct {
	AllowedNetworkTargets []NetworkTargetRule `json:"allowedNetworkTargets"`
	FsWriteAccess         map[string]string   `json:"fsWriteAccess"`
	StdoutAccess          bool                `json:"stdoutAccess"`
	StderrAccess          bool                `json:"stderrAccess"`
}

type NetworkTargetRule struct {
	Pattern string `json:"pattern"`
}
