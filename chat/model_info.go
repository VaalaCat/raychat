package chat

import (
	"github.com/imroc/req/v3"
)

func (r *RayChat) GetSupportedModels() map[string]string {
	c := req.C().SetCommonHeaders(map[string]string{
		"Accept":          "application/json",
		"Accept-Language": "zh-CN,zh-Hans;q=0.9",
		"User-Agent":      "Raycast/0 CFNetwork/1408.0.4 Darwin/22.5.0",
		"Content-Type":    "application/json",
		"Authorization":   "Bearer " + r.Token,
	})

	resp := GetAIInfoResponse{}

	res, err := c.R().SetSuccessResult(&resp).Get("https://backend.raycast.com/api/v1/ai/models")
	if err != nil {
		Logger().WithError(err).Panic("get model info failed")
	}
	if res.StatusCode != 200 {
		Logger().WithField("status code", res.StatusCode).Panic("get model info failed")
	}
	Logger().Infof("get model info success, support those models: [%+v], resp: [%+v]", resp.SupporedModels(), resp)

	return resp.SupporedModels()
}
