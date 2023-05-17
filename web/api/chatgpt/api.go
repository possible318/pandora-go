package chatgpt

import (
	"bufio"
	"encoding/json"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	api "github.com/pandora_go/server/web"
	reqTypes "github.com/pandora_go/web/typings/req"
	"github.com/spf13/viper"
	"strings"
)

// ListModels 获取模型列表
func ListModels(c *gin.Context) {
	res, err := api.ListModels()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}

// ListConversations 获取会话列表
func ListConversations(c *gin.Context) {
	offset, ok := c.GetQuery("offset")
	if !ok {
		offset = "0"
	}
	limit, ok := c.GetQuery("limit")
	if !ok {
		limit = "20"
	}
	res, err := api.ListConversations(offset, limit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}

// GetConversation 获取会话内容
func GetConversation(c *gin.Context) {
	id := c.Param("id")

	res, err := api.GetConversation(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}

// ClearConversations 清空会话列表
func ClearConversations(c *gin.Context) {
	res, err := api.ClearConversations()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}

// DelConversation  删除会话
func DelConversation(c *gin.Context) {
	var params reqTypes.ConversationRequest
	c.ShouldBindJSON(&params)
	if params.Title != "" {
		params.IsVisible = false
	}
	id := c.Param("id")

	res, err := api.DelConversation(id, params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}

// SetConversationTitle  重命名
func SetConversationTitle(c *gin.Context) {
	var params reqTypes.ConversationRequest
	c.ShouldBindJSON(&params)
	if params.Title != "" {
		params.IsVisible = true
	}
	id := c.Param("id")

	res, err := api.SetConversationTitle(id, params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})

	}
	c.JSON(http.StatusOK, res)

}

// GenConversationTitle 生成会话标题
func GenConversationTitle(c *gin.Context) {
	var params reqTypes.GenerateTitleRequest
	c.ShouldBindJSON(&params)
	id := c.Param("id")

	res, err := api.GenConversationTitle(id, params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorMessage": err.Error()})

	}
	c.JSON(http.StatusOK, res)
}

func Session(c *gin.Context) {
	ret := map[string]interface{}{
		"user": map[string]interface{}{
			"id":      "user-000000000000000000000000",
			"name":    "admin@openai.com",
			"email":   "admin@openai.com",
			"image":   "",
			"picture": "",
			"groups":  []string{""},
		},
		"expires":     "2089-08-08T23:59:59.999Z",
		"accessToken": "secret",
	}
	c.JSON(http.StatusOK, ret)
}

func Check(c *gin.Context) {
	ret := map[string]interface{}{
		"account_plan": map[string]interface{}{
			"is_paid_subscription_active":       true,
			"subscription_plan":                 "chatgptplusplan",
			"account_user_role":                 "account-owner",
			"was_paid_customer":                 true,
			"has_customer_object":               true,
			"subscription_expires_at_timestamp": 3774355199,
		},
		"user_country": "US",
		"features": []string{
			"model_switcher",
			"dfw_message_feedback",
			"dfw_inline_message_regen_comparison",
			"model_preview",
			"system_message",
			"can_continue",
		},
	}
	c.JSON(http.StatusOK, ret)
}

// Chat 首页
func Chat(c *gin.Context) {
	id := c.Param("id")
	query := map[string]any{}
	if id != "" {
		query = map[string]any{"chatId": id}
	}
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"pandora_base":   c.Request.URL.Host,
		"query":          query,
		"pandora_sentry": viper.Get("sentry"),
	})
}

func GetTokenKey(c *gin.Context) string {
	// 请求头拿数据
	tokenKey := c.Request.Header.Get("X-Use-Token")
	if tokenKey != "" {
		// 从cookie中拿
		tokenKey, _ = c.Cookie("token-key")
	}
	return tokenKey
}

func SetCookie(c *gin.Context, tokenKey string) {
	// 设置cookie
	c.SetCookie("token-key", tokenKey, 60*60*24*30, "/", "", false, true)
}

// Talk 对话
func Talk(c *gin.Context) {
	var params reqTypes.ApiRequest
	c.ShouldBindJSON(&params)
	params.Stream = true
	c.Writer.Header().Set("Content-Type", "text/event-stream")

	chatParam := map[string]interface{}{
		"action": "next",
		"messages": []map[string]interface{}{
			{
				"id":   params.MessageId,
				"role": "user",
				"author": map[string]interface{}{
					"role": "user",
				},
				"content": map[string]interface{}{
					"content_type": "text",
					"parts":        []string{params.Prompt},
				},
			},
		},
		"model":             params.Model,
		"parent_message_id": params.ParentMessageId,
	}

	if params.ConversationId != "" {
		chatParam["conversation_id"] = params.ConversationId
	}

	response, err := api.Conversation(chatParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			c.JSON(response.StatusCode, err)
		}
		c.JSON(response.StatusCode, gin.H{
			"error":   "error sending request",
			"content": errorResponse,
		})
		return
	}
	HandleSseRes(c, response)
}

// Goon 继续对话
func Goon(c *gin.Context) {
	var params reqTypes.GoonRequest
	c.ShouldBindJSON(&params)

	params.Stream = true

	newParam := map[string]interface{}{
		"action":            "continue",
		"conversation_id":   params.ConversationId,
		"model":             params.Model,
		"parent_message_id": params.ParentMessageId,
	}

	response, err := api.Conversation(newParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			c.JSON(response.StatusCode, err)
		}
		c.JSON(response.StatusCode, gin.H{
			"error":   "error sending request",
			"content": errorResponse,
		})
		return
	}
	HandleSseRes(c, response)
}

// Regenerate 重新生成
func Regenerate(c *gin.Context) {
	var params reqTypes.ApiRequest
	c.ShouldBindJSON(&params)
	params.Stream = true
	c.Header("Content-Type", "text/event-stream")

	chatParam := map[string]interface{}{
		"action": "variant",
		"messages": []map[string]interface{}{
			{
				"id":   params.MessageId,
				"role": "user",
				"author": map[string]interface{}{
					"role": "user",
				},
				"content": map[string]interface{}{
					"content_type": "text",
					"parts":        []string{params.Prompt},
				},
			},
		},
		"model":             params.Model,
		"parent_message_id": params.ParentMessageId,
		"conversation_id":   params.ConversationId,
	}

	response, err := api.Conversation(chatParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		var errorResponse map[string]interface{}
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			c.JSON(response.StatusCode, err)
		}
		c.JSON(response.StatusCode, gin.H{
			"error":   "error sending request",
			"content": errorResponse,
		})
		return
	}
	HandleSseRes(c, response)
}

func HandleSseRes(c *gin.Context, resp *http.Response) {
	reader := bufio.NewReader(resp.Body)
	for {
		if c.Request.Context().Err() != nil {
			break
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "event") ||
			strings.HasPrefix(line, "data: 20") ||
			line == "" {
			continue
		}
		c.Writer.Write([]byte(line + "\n\n"))
		c.Writer.Flush()
	}
}
