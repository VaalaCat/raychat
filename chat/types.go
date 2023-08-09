package chat

import (
	"encoding/json"
	"raychat/auth"
	"strings"
	"time"

	"github.com/samber/lo"
)

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Stream      bool            `json:"stream"`
	Temperature float64         `json:"temperature"`
}

func (r OpenAIRequest) ToRayChatRequest(a *auth.RaycastAuth) RayChatRequest {
	messages := make([]RayChatMessage, len(r.Messages))
	for i, m := range r.Messages {
		messages[i] = m.ToRayChatMessage()
	}
	if r.Temperature == 0 {
		r.Temperature = 1
	}

	return RayChatRequest{
		Debug:             false,
		Locale:            "en-CN",
		Provider:          "openai",
		Model:             r.GetRequestModel(a),
		Temperature:       r.Temperature,
		SystemInstruction: "markdown",
		Messages:          messages,
	}
}

func (r OpenAIRequest) GetRequestModel(a *auth.RaycastAuth) string {
	model := r.Model
	supporedModels := []string{}
	for _, m := range a.LoginResp.User.AiChatModels {
		supporedModels = append(supporedModels, m.Model)
	}
	if !lo.Contains(supporedModels, r.Model) {
		model = "gpt-3.5-turbo"
	}
	return model
}

type RayChatRequest struct {
	Debug             bool             `json:"debug"`
	Locale            string           `json:"locale"`
	Messages          []RayChatMessage `json:"messages"`
	Provider          string           `json:"provider"`
	Model             string           `json:"model"`
	Temperature       float64          `json:"temperature"`
	SystemInstruction string           `json:"system_instruction"`
}

type Content struct {
	Text string `json:"text"`
}

type RayChatMessage struct {
	Content Content `json:"content"`
	Author  string  `json:"author"`
}

func (m RayChatMessage) ToOpenAIMessage() OpenAIMessage {
	return OpenAIMessage{
		Role:    m.Author,
		Content: m.Content.Text,
	}
}

type RayChatStreamResponse struct {
	Text         string  `json:"text"`
	FinishReason *string `json:"finish_reason"`
}

func (r RayChatStreamResponse) FromEventString(origin string) RayChatStreamResponse {
	selection := strings.Replace(origin, "data: ", "", 1)
	if len(selection) == 0 {
		return RayChatStreamResponse{}
	}
	err := json.Unmarshal([]byte(selection), &r)
	if err != nil {
		panic(err)
	}
	return r
}

func (r RayChatStreamResponse) ToOpenAISteamResponse(model string) OpenAIStreamResponse {

	resp := OpenAIStreamResponse{
		ID:      "chatcmpl-" + generateRandomString(29),
		Object:  "chat.completion.chunk",
		Created: int(time.Now().Unix()),
		Model:   model,
		Choices: []StreamChoices{
			{
				Index:        0,
				FinishReason: r.FinishReason,
			},
		},
	}
	if len(r.Text) != 0 {
		resp.Choices[0].Delta = Delta{
			Content: r.Text,
		}
	}

	return resp
}

type RayChatStreamResponses []RayChatStreamResponse

func (r RayChatStreamResponses) ToOpenAIResponse(model string) OpenAIResponse {
	content := ""
	for _, resp := range r {
		content += resp.Text
	}
	return OpenAIResponse{
		ID:      "chatcmpl-" + generateRandomString(29),
		Object:  "chat.completion",
		Created: int(time.Now().Unix()),
		Choices: []Choices{
			{
				Index: 0,
				Message: OpenAIMessage{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: lo.ToPtr("stop"),
			},
		},
		Usage: Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
	}
}

type OpenAIResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Choices []Choices `json:"choices"`
	Usage   Usage     `json:"usage"`
}

func (o OpenAIResponse) ToEventString() string {
	bytesRsp, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return "data: " + string(bytesRsp)
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (m OpenAIMessage) ToRayChatMessage() RayChatMessage {
	role := m.Role
	if m.Role == "system" {
		role = "user"
	}
	return RayChatMessage{
		Content: Content{
			Text: m.Content,
		},
		Author: role,
	}
}

type Choices struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason *string       `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAIStreamResponse struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int             `json:"created"`
	Model   string          `json:"model"`
	Choices []StreamChoices `json:"choices"`
}

func (o OpenAIStreamResponse) ToEventString() string {
	bytesRsp, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return "data: " + string(bytesRsp)
}

type Delta struct {
	Content string `json:"content,omitempty"`
}

type StreamChoices struct {
	Index        int     `json:"index"`
	Delta        Delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}
