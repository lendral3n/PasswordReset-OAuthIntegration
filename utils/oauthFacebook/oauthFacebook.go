package oauthfacebook

import (
	"context"
	"emailnotifl3n/app/config"
	"emailnotifl3n/features/user"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type FacebookOauth struct {
	oauthConfig *oauth2.Config
}

type FacebookOauthToken struct {
	Access_token string
}

type FacebookInterface interface {
	GetAuthURL() string
	GetFacebookOauthToken(code string) (*FacebookOauthToken, error)
	GetFacebookUser(access_token string) (*user.Core, error)
}

func New() FacebookInterface {
	return &FacebookOauth{
		oauthConfig: &oauth2.Config{
			RedirectURL:  config.FB_URL,
			ClientID:     config.CLIENT_ID_FB,
			ClientSecret: config.CLIENT_SECRET_FB,
			Scopes:       config.SCOPES_FB,
			Endpoint:     facebook.Endpoint,
		},
	}
}

// GetAuthURL implements FacebookInterface.
func (facebook *FacebookOauth) GetAuthURL() string {
	return facebook.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// GetFacebookOauthToken implements FacebookInterface.
func (facebook *FacebookOauth) GetFacebookOauthToken(code string) (*FacebookOauthToken, error) {
	token, err := facebook.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	return &FacebookOauthToken{Access_token: token.AccessToken}, nil

}

// GetFacebookUser implements FacebookInterface.
func (facebook *FacebookOauth) GetFacebookUser(access_token string) (*user.Core, error) {
	response, err := http.Get("https://graph.facebook.com/me?fields=id,name,email,picture&access_token=" + url.QueryEscape(access_token))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info from Facebook")
	}

	var fbUserRes map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&fbUserRes)
	if err != nil {
		return nil, err
	}

	userBody := &user.Core{
		Verified:         true,
		RegistrationType: "Facebook",
	}

	if email, ok := fbUserRes["email"].(string); ok {
		userBody.Email = email
	}

	if name, ok := fbUserRes["name"].(string); ok {
		userBody.Name = name
	}

	if picture, ok := fbUserRes["picture"].(map[string]interface{}); ok {
		if data, ok := picture["data"].(map[string]interface{}); ok {
			if url, ok := data["url"].(string); ok {
				userBody.PhotoProfile = url
			}
		}
	}
	return userBody, nil
}
