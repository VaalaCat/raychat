package auth

type StepOneResponse struct {
	AuthenticityToken string `json:"authenticity_token"`
}

type LoginResponse struct {
	User              User   `json:"user"`
	RedirectTo        string `json:"redirect_to"`
	AuthenticityToken string `json:"authenticity_token"`
}

type AiChatModels struct {
	Model    string `json:"model"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

type User struct {
	AiChatModels           []AiChatModels `json:"ai_chat_models"`
	EligibleForGpt4        bool           `json:"eligible_for_gpt4"`
	Admin                  bool           `json:"admin"`
	HasActiveSubscription  bool           `json:"has_active_subscription"`
	HasRunningSubscription bool           `json:"has_running_subscription"`
	Email                  string         `json:"email"`
	Name                   string         `json:"name"`
	Handle                 string         `json:"handle"`
	Username               string         `json:"username"`
}

type StepFiveResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	CreatedAt   int    `json:"created_at"`
	Data        Data   `json:"data"`
}

type Data struct {
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Name     string      `json:"name"`
	Avatar   interface{} `json:"avatar"`
}
