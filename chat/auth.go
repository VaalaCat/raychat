package chat

import (
	"raychat/auth"
	"raychat/settings"
)

var token string

func init() {
	if settings.Get().Token != "" {
		token = settings.Get().Token
		return
	}
	token = (&auth.RaycastAuth{
		ClientID:     settings.Get().ClientID,
		ClientSecret: settings.Get().ClientSecret,
		Email:        settings.Get().Email,
		Password:     settings.Get().Password,
	}).Login()
}

func getToken() string {
	return token
}
