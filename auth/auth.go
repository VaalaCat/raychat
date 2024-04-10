package auth

import (
	"net/url"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Entry {
	return logrus.WithField("prefix", "raycast")
}

type RaycastAuth struct {
	ClientID     string
	ClientSecret string
	Email        string
	Password     string
	LoginResp    LoginResponse
}

func (r *RaycastAuth) Login() string {
	cli := req.C().
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5.2 Safari/605.1.15")
	r1 := r.stepOne(cli)
	r.stepTwo(cli, r1, r.ClientID)
	r3 := r.stepThree(cli, r.Email, r.Password)
	r4 := r.stepFour(cli, r3.RedirectTo)
	r5 := r.stepFive(r4, r.ClientID, r.ClientSecret)
	return r5.AccessToken
}

func (r *RaycastAuth) stepOne(c *req.Client) StepOneResponse {
	var resp StepOneResponse
	rawResp, err := c.R().SetSuccessResult(&resp).Get("https://www.raycast.com/frontend_api/session")
	if err != nil {
		Logger().WithError(err).Panic("step one failed, raw response: ", rawResp)
	}
	Logger().Info("step one success, authenticity token: ", resp.AuthenticityToken)
	return resp
}

func (r *RaycastAuth) stepTwo(c *req.Client, prev StepOneResponse, clientID string) {
	c.R().
		SetHeaders(map[string]string{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "zh-CN,zh-Hans;q=0.9",
			"Sec-Fetch-Mode":  "navigate",
			"Sec-Fetch-Site":  "none",
			"Sec-Fetch-Dest":  "document",
			"Accept-Encoding": "gzip, deflate, br",
		}).
		Get("https://www.raycast.com/oauth/authorize" +
			"?" + "client_id=" + clientID +
			"&" + "redirect_uri=https://raycast.com/redirect?packageName%3DRaycast%2520Account" +
			"&" + "state=%7B%22id%22:%2268D61350-A340-4E2C-8924-130867700072%22,%22flavor%22:%22release%22%7D" +
			"&" + "response_type=code" +
			"&" + "audience=" +
			"&" + "scope=")
}

func (r *RaycastAuth) stepThree(c *req.Client, email, password string) LoginResponse {
	var resp LoginResponse
	rawResp, err := c.R().SetSuccessResult(&resp).
		SetHeaders(map[string]string{
			"Accept":          "application/json",
			"Accept-Language": "zh-CN,zh-Hans;q=0.9",
			"Sec-Fetch-Site":  "same-origin",
			"Sec-Fetch-Mode":  "cors",
			"Sec-Fetch-Dest":  "empty",
			"Content-Type":    "application/json",
			"Origin":          "https://www.raycast.com",
			"Referer":         "https://www.raycast.com/users/sign_in",
			"X-CSRF-Token":    getCSRFToken(c),
		}).
		SetBody(map[string]map[string]string{
			"user": {
				"email":    email,
				"password": password,
			},
		}).
		Post("https://www.raycast.com/frontend_api/session")
	if err != nil {
		Logger().WithError(err).Panic("login failed, raw response: ", rawResp.String())
	}
	Logger().Infof("login success, resp: %+v", rawResp.String())
	r.LoginResp = resp
	return resp
}

func (r *RaycastAuth) stepFour(c *req.Client, redirUrl string) string {
	url := "https://www.raycast.com" + redirUrl
	Logger().Info("redirect url: ", url)
	resp, err := c.SetRedirectPolicy(req.NoRedirectPolicy()).R().SetHeaders(map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "zh-CN,zh-Hans;q=0.9",
		"Sec-Fetch-Site":  "same-origin",
		"Sec-Fetch-Mode":  "navigate",
		"Sec-Fetch-Dest":  "document",
		"Content-Type":    "application/json",
		"Referer":         "https://www.raycast.com/users/sign_in",
		"X-CSRF-Token":    getCSRFToken(c),
	}).Get(url)
	if err != nil {
		Logger().WithError(err).Panic("step four redirect failed")
	}
	redir := resp.GetHeader("Location")
	return redir
}

func (r *RaycastAuth) stepFive(redirUrl, clientID, clientSecret string) StepFiveResponse {
	Logger().Info("redirect url: ", redirUrl)
	parsedURL, err := url.Parse(redirUrl)
	if err != nil {
		Logger().WithError(err).Panic("parse redirect url failed")
		return StepFiveResponse{}
	}
	qp := map[string]string{}
	queryParams := parsedURL.Query()
	for key, values := range queryParams {
		if len(values) > 0 {
			qp[key] = values[0]
		}
	}

	var resp StepFiveResponse
	cli := req.C().SetUserAgent("Raycast/0 CFNetwork/1408.0.4 Darwin/22.5.0")
	rawResp, err := cli.R().SetSuccessResult(&resp).
		SetHeaders(map[string]string{
			"Content-Type":    "application/x-www-form-urlencoded",
			"Accept":          "application/json",
			"Accept-Language": "zh-CN,zh-Hans;q=0.9",
		}).
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     clientID,
			"code":          qp["code"],
			"redirect_uri":  "https://raycast.com/redirect?packageName=Raycast%20Account",
			"client_secret": clientSecret,
		}).
		Post("https://www.raycast.com/oauth/token")
	if err != nil {
		Logger().WithError(err).Panic("step five failed, raw response: ", rawResp.String())
	}
	Logger().Info("step five success, resp: ", resp)
	return resp
}

func getCSRFToken(c *req.Client) string {
	cookies, err := c.GetCookies("https://www.raycast.com")
	if err != nil {
		Logger().WithError(err).Panic("get csrf token failed")
	}
	for _, t := range cookies {
		if t.Name == "csrf_token" {
			return t.Value
		}
	}
	return ""
}
