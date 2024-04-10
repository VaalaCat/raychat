package chat

import (
	"bufio"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ChatEndpoint(c *gin.Context) {
	originReq := &OpenAIRequest{}
	if err := c.Copy().ShouldBindJSON(originReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := Cli(getToken()).Chat(originReq.ToRayChatRequest(getAuthInstance()))
	if err != nil || r.StatusCode != http.StatusOK {
		data, err := io.ReadAll(r.Body)
		logrus.WithError(err).Errorf("request to raycast error, request: %+v, body: %+v", r, string(data))
		c.JSON(http.StatusBadRequest, gin.H{"error": "request to raycast error", "code": 400})
		return
	}
	switch originReq.Stream {
	case true:
		streamResp(c, originReq, r)
	default:
		plainResp(c, originReq, r)
	}
}

func plainResp(c *gin.Context, req *OpenAIRequest, resp *http.Response) {
	defer resp.Body.Close()

	model, _ := req.GetRequestModel(getAuthInstance())

	rayChatResps := *new(RayChatStreamResponses)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		event := scanner.Text()
		if len(event) == 0 {
			continue
		}
		rayChatResp := RayChatStreamResponse{}.FromEventString(event)
		rayChatResps = append(rayChatResps, rayChatResp)
	}
	if scanner.Err() != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": scanner.Err().Error()})
		return
	}
	openaiResp := rayChatResps.ToOpenAIResponse(model)
	c.JSON(http.StatusOK, openaiResp)
}

func streamResp(c *gin.Context, req *OpenAIRequest, resp *http.Response) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	_, ok := c.Writer.(http.Flusher)
	if !ok {
		logrus.Panic("server not support")
	}
	defer func() {
		c.Writer.WriteString("data: [DONE]\n\n")
		c.Writer.Flush()
		resp.Body.Close()
	}()

	model, _ := req.GetRequestModel(getAuthInstance())

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		event := scanner.Text()
		if len(event) == 0 {
			c.Writer.WriteString("\n")
			c.Writer.Flush()
			continue
		}
		rayChatResp := RayChatStreamResponse{}.FromEventString(event)
		openAIResp := rayChatResp.ToOpenAISteamResponse(model)
		eventResp := openAIResp.ToEventString()
		_, err := c.Writer.WriteString(eventResp + "\n")
		if err != nil {
			c.Writer.WriteString("data: {\"finish_reason\":\"stop\"}" + "\n")
			return
		}
		c.Writer.Flush()
	}
	if scanner.Err() != nil {
		return
	}
}
