package chat

import (
	"raychat/auth"
	"raychat/settings"
)

var (
	authInstance *auth.RaycastAuth
	token        string
)

func init() {
	Init()
}

func Init() {
	if settings.Get().Token != "" {
		token = settings.Get().Token
		return
	}
	authInstance = &auth.RaycastAuth{
		ClientID:     settings.Get().ClientID,
		ClientSecret: settings.Get().ClientSecret,
		Email:        settings.Get().Email,
		Password:     settings.Get().Password,
	}
	token = authInstance.Login()
}

func getToken() string {
	return token
}

func getAuthInstance() *auth.RaycastAuth {
	return authInstance
}
