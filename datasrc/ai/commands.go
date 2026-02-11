package ai

const (
	CMD_SEND_PROMPT = "sendPrompt"
	CMD_GET_MODELS  = "getModels"
)

func AllCommands() map[string]any {
	return map[string]any{
		"CMD_SEND_PROMPT": CMD_SEND_PROMPT,
		"CMD_GET_MODELS":  CMD_GET_MODELS,
	}
}
