package ai

import "fmt"

// Message represents a single message in a conversation.
// Role should be "user" or "assistant". This is vendor-agnostic;
// each datasource implementation maps these roles to their SDK's conventions.
type Message struct {
	Role    string `json:"role" mapstructure:"role"`
	Content string `json:"content" mapstructure:"content"`
}

// SendPromptParams contains parameters for the sendPrompt command.
// Messages contains the conversation history. The caller decides whether to
// send the full history (multi-turn) or a single message (stateless).
type SendPromptParams struct {
	Messages          []Message `json:"messages" mapstructure:"messages" validate:"required"`
	Model             string    `json:"model" mapstructure:"model" validate:"required"`
	SystemInstruction string    `json:"systemInstruction,omitempty" mapstructure:"systemInstruction"`
	Temperature       float32   `json:"temperature,omitempty" mapstructure:"temperature"`
	MaxOutputTokens   int32     `json:"maxOutputTokens,omitempty" mapstructure:"maxOutputTokens"`
}

func (p SendPromptParams) Validate() error {
	if len(p.Messages) == 0 {
		return fmt.Errorf("messages is required")
	}
	if p.Model == "" {
		return fmt.Errorf("model is required")
	}
	return nil
}

// SendPromptResponse is the response type for the sendPrompt command.
type SendPromptResponse struct {
	Text  string        `json:"text"`
	Model string        `json:"model"`
	Usage UsageMetadata `json:"usage"`
}

// UsageMetadata contains token usage information from the AI provider.
type UsageMetadata struct {
	PromptTokens   int32 `json:"promptTokens"`
	ResponseTokens int32 `json:"responseTokens"`
	TotalTokens    int32 `json:"totalTokens"`
}

// ModelInfo describes an available AI model.
type ModelInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

// GetModelsResponse is the response type for the getModels command.
type GetModelsResponse struct {
	Models []ModelInfo `json:"models"`
}
