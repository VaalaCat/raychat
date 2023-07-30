package chat

import (
	"bufio"
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
	if err != nil {
		panic(err)
	}
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
		r.Body.Close()
	}()

	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		event := scanner.Text()
		if len(event) == 0 {
			c.Writer.WriteString("\n")
			c.Writer.Flush()
			continue
		}
		rayChatResp := RayChatStreamResponse{}.FromEventString(event)
		openAIResp := rayChatResp.ToOpenAISteamResponse(originReq.GetRequestModel(getAuthInstance()))
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
