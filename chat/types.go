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
	messages := make([]RayChatMessage, 0, len(r.Messages))
	for _, m := range r.Messages {
		if m.Role == "system" {
			continue
		}
		messages = append(messages, m.ToRayChatMessage())
	}
	if r.Temperature == 0 {
		r.Temperature = 1
	}

	model, provider := r.GetRequestModel(a)

	resp := RayChatRequest{
		Debug:             false,
		Locale:            "en-CN",
		Provider:          provider,
		Model:             model,
		Temperature:       r.Temperature,
		SystemInstruction: "markdown",
		Messages:          messages,
	}

	additionalSystemInstructions := r.GetSystemMessage().Content
	if additionalSystemInstructions != "" {
		resp.AdditionalSystemInstructions = additionalSystemInstructions
	}

	return resp
}

func (r OpenAIRequest) GetRequestModel(a *auth.RaycastAuth) (string, string) {
	model := r.Model
	supporedModels := lo.Keys(models)
	for _, m := range a.LoginResp.User.AiChatModels {
		supporedModels = append(supporedModels, m.Model)
	}
	if a.LoginResp.User.EligibleForGpt4 {
		supporedModels = append(supporedModels, "gpt-4")
	}

	if !lo.Contains(supporedModels, r.Model) {
		model = "gpt-3.5-turbo"
	}
	return model, models[model]
}

func (r OpenAIRequest) GetSystemMessage() OpenAIMessage {
	additionalSystem := ""
	for _, m := range r.Messages {
		if m.Role == "system" {
			additionalSystem += m.Content
		}
	}
	return OpenAIMessage{
		Role:    "system",
		Content: additionalSystem,
	}
}

func (r OpenAIRequest) GetNoneSystemMessage() []OpenAIMessage {
	msgs := []OpenAIMessage{}
	for _, m := range r.Messages {
		if m.Role != "system" {
			msgs = append(msgs, m)
		}
	}
	return msgs
}

type RayChatRequest struct {
	Debug                        bool             `json:"debug"`
	Locale                       string           `json:"locale"`
	Messages                     []RayChatMessage `json:"messages"`
	Source                       string           `json:"source"`
	Provider                     string           `json:"provider"`
	Model                        string           `json:"model"`
	Temperature                  float64          `json:"temperature"`
	SystemInstruction            string           `json:"system_instruction"`
	AdditionalSystemInstructions string           `json:"additional_system_instructions,omitempty"`
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
	Text         string      `json:"text"`
	FinishReason *string     `json:"finish_reason"`
	Err          interface{} `json:"error"`
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
	if strings.Contains(selection, "error") {
		Logger().WithError(err).Errorf("request to raycast error, body: %+v", origin)
		return RayChatStreamResponse{}
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

type GetAIInfoResponse struct {
	Models        []ModelInfo `json:"models"`
	DefaultModels struct {
		Chat        string `json:"chat"`
		QuickAi     string `json:"quick_ai"`
		Commands    string `json:"commands"`
		API         string `json:"api"`
		EmojiSearch string `json:"emoji_search"`
	} `json:"default_models"`
}

func (m GetAIInfoResponse) SupporedModels() map[string]string {
	models := map[string]string{}
	for _, model := range m.Models {
		models[model.Model] = model.Provider
	}
	return models
}

type ModelInfo struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	Status                 any      `json:"status"`
	Features               []string `json:"features"`
	Suggestions            []any    `json:"suggestions"`
	InBetterAiSubscription bool     `json:"in_better_ai_subscription"`
	Model                  string   `json:"model"`
	Provider               string   `json:"provider"`
	ProviderName           string   `json:"provider_name"`
	ProviderBrand          string   `json:"provider_brand"`
	Speed                  int      `json:"speed"`
	Intelligence           int      `json:"intelligence"`
	RequiresBetterAi       bool     `json:"requires_better_ai"`
	Context                int      `json:"context"`
	Capabilities           struct {
		WebSearch       string `json:"web_search,omitempty"`
		ImageGeneration string `json:"image_generation,omitempty"`
	} `json:"capabilities,omitempty"`
}
