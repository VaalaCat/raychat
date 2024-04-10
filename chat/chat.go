package chat

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	url = "https://backend.raycast.com/api/v1/ai/chat_completions"
)

type RayChat struct {
	Token string
}

func Cli(token string) *RayChat {
	return &RayChat{
		Token: token,
	}
}

func (r *RayChat) Chat(request RayChatRequest) (*http.Response, error) {
	rawReq, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	payload := bytes.NewReader(rawReq)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "zh-CN,zh-Hans;q=0.9")
	req.Header.Add("User-Agent", "Raycast/0 CFNetwork/1408.0.4 Darwin/22.5.0")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+r.Token)

	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return res, err
	}

	return res, nil
}