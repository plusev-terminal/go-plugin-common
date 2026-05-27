package ai

const (
	CMD_SEND_PROMPT        = "sendPrompt"
	CMD_SEND_PROMPT_STREAM = "sendPromptStream"
	CMD_GET_MODELS         = "getModels"
)

func AllCommands() map[string]any {
	return map[string]any{
		"CMD_SEND_PROMPT":        CMD_SEND_PROMPT,
		"CMD_SEND_PROMPT_STREAM": CMD_SEND_PROMPT_STREAM,
		"CMD_GET_MODELS":         CMD_GET_MODELS,
	}
}
