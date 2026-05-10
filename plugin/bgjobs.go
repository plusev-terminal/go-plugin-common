package plugin

// PluginBgJobCreate is sent by a plugin to create a new background job.
type PluginBgJobCreate struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Icon        string              `json:"icon,omitempty"`
	Description string              `json:"description,omitempty"`
	Progress    float64             `json:"progress"`
	Actions     []PluginBgJobAction `json:"actions,omitempty"`
}

// PluginBgJobAction represents an action button on a background job.
type PluginBgJobAction struct {
	Label   string `json:"label"`
	Action  string `json:"action"`
	Data    string `json:"data,omitempty"`
	Variant string `json:"variant,omitempty"`
}

// PluginBgJobUpdate is sent by a plugin to update a running background job.
type PluginBgJobUpdate struct {
	ID          string   `json:"id"`
	Description string   `json:"description,omitempty"`
	Progress    *float64 `json:"progress,omitempty"`
}

// PluginBgJobComplete is sent by a plugin to mark a background job as completed.
type PluginBgJobComplete struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
}

// PluginBgJobFail is sent by a plugin to mark a background job as failed.
type PluginBgJobFail struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
}

// PluginBgJobResult is returned to the plugin after a background job operation.
type PluginBgJobResult struct {
	Result bool   `json:"result"`
	ID     string `json:"id,omitempty"`
	Error  string `json:"error,omitempty"`
}
