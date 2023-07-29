package chat

import (
	"encoding/json"
	"strings"
	"time"
)

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Stream      bool            `json:"stream"`
	Temperature float64         `json:"temperature"`
}

func (r OpenAIRequest) ToRayChatRequest() RayChatRequest {
	messages := make([]Message, len(r.Messages))
	for i, m := range r.Messages {
		messages[i] = m.ToMessage()
	}
	if r.Temperature == 0 {
		r.Temperature = 1
	}
	return RayChatRequest{
		Debug:             false,
		Locale:            "en-CN",
		Provider:          "openai",
		Model:             r.Model,
		Temperature:       r.Temperature,
		SystemInstruction: "markdown",
		Messages:          messages,
	}
}

type RayChatRequest struct {
	Debug             bool      `json:"debug"`
	Locale            string    `json:"locale"`
	Messages          []Message `json:"messages"`
	Provider          string    `json:"provider"`
	Model             string    `json:"model"`
	Temperature       float64   `json:"temperature"`
	SystemInstruction string    `json:"system_instruction"`
}

type Content struct {
	Text string `json:"text"`
}

type Message struct {
	Content Content `json:"content"`
	Author  string  `json:"author"`
}

func (m Message) ToOpenAIMessage() OpenAIMessage {
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

func (r RayChatStreamResponse) ToOpenAISteamResponse() OpenAIStreamResponse {

	resp := OpenAIStreamResponse{
		ID:      "chatcmpl-" + generateRandomString(29),
		Object:  "chat.completion.chunk",
		Created: int(time.Now().Unix()),
		Model:   "gpt-3.5-turbo-0613",
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

func (m OpenAIMessage) ToMessage() Message {
	role := m.Role
	if m.Role == "system" {
		role = "assistant"
	}
	return Message{
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
