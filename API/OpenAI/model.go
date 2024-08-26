package openai

import (
	"github.com/otiai10/openaigo"
)

const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
	ChatMessageRoleFunction  = "function"
	ChatMessageRoleTool      = "tool"
)

type QueuePrompt struct {
	Emotion     string
	CallBack    chan ChatResult
	PromptGroup openaigo.Message
}

type SavePrompt struct {
	TokenUse    int
	PromptGroup openaigo.Message
}

type ChatResult struct {
	Message string `json:"Message"`
	Emotion string `json:"Emotion"`
}
